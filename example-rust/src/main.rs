use glob::glob;
use itertools::Itertools;
use lexical_sort::{natural_lexical_cmp, StringSort};
use once_cell::sync::Lazy;
use regex::Regex;
use std::path::Path;
use tokio::fs;
use yaml_rust::YamlLoader;

const SORT: bool = true;
const OUTDIR: &str = "words";
const FILE_DIR: &str = "../data/??/**/*.yml";
static SEPARATOR_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"[\W\d]+").unwrap());

async fn create_outdir() -> tokio::io::Result<()> {
    fs::create_dir_all(OUTDIR).await
}

async fn read_file(path: &Path) -> String {
    let raw = fs::read_to_string(path).await.unwrap();
    raw.to_lowercase().replace('\n', " ")
}

fn get_unique_token(src: &str) -> Vec<&str> {
    let mut data = SEPARATOR_REGEX.split(src).unique().collect::<Vec<_>>();

    if SORT {
        data.string_sort_unstable(natural_lexical_cmp);
    }

    data
}

async fn get_filename_from_meta(path: &Path) -> anyhow::Result<String> {
    let yaml = fs::read_to_string(path).await?;
    let docs = YamlLoader::load_from_str(&yaml)?;
    let meta = &docs[0];

    let lang = meta["lang"]
        .as_str()
        .ok_or_else(|| anyhow::anyhow!("code not found"))?;

    let code = meta["code"]
        .as_str()
        .ok_or_else(|| anyhow::anyhow!("code not found"))?;

    Ok(format!("{}/{}-{}.txt", OUTDIR, lang, code))
}

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let start = std::time::Instant::now();
    let path = glob(FILE_DIR).expect("failed to read glob pattern");

    let submissions = path.map(|entry| {
        tokio::spawn(async {
            let yaml_path = entry.expect("should be existed");
            let txt_path = yaml_path.with_extension("txt");

            let outdir_submission =
                tokio::spawn(async { create_outdir().await.expect("unable to create outdir") });

            let read_text_file_submission = tokio::spawn(async move {
                let data = read_file(&txt_path).await;
                let tokens = get_unique_token(&data);

                tokens.join("\n")
            });

            let filename_submission = tokio::spawn(async move {
                get_filename_from_meta(&yaml_path)
                    .await
                    .expect("should be existed")
            });

            let (tokens, filename, _) = tokio::join!(
                read_text_file_submission,
                filename_submission,
                outdir_submission,
            );

            fs::write(
                filename.expect("failed to run filename"),
                tokens.expect("failed to get tokens"),
            )
            .await
            .expect("failed to write");
        })
    });

    futures::future::join_all(submissions).await;

    println!("{:?}", start.elapsed());
    Ok(())
}