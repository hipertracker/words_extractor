import glob
from icu import Collator, Locale
import os
import re
import shutil
import time
import yaml


def worker(path, collator, separator, outdir, with_sorting):
    filepath = path.replace(".yml", ".txt")
    print(f"Processing {filepath}")
    with open(filepath) as file:
        text = file.read().lower().rstrip()
        words = set(re.split(separator, text))
    with open(path) as file:
        meta = yaml.safe_load(file)
    with open(f"{outdir}/extracted-words-for-{meta['label']}.txt", "w") as file:
        if with_sorting:
            words = sorted(words, key=collator.getSortKey)
        file.write("\n".join(words))


if __name__ == "__main__":
    t = time.time()

    outdir = "words"
    if os.path.exists(outdir):
        shutil.rmtree(outdir)
    os.makedirs(outdir)

    collator = Collator.createInstance(Locale("pl_PL.UTF-8"))
    separator = re.compile("[\W\d]+")
    for path in glob.glob("../data/pl/**/*.yml", recursive=True):
        worker(
            path=path,
            collator=collator,
            separator=separator,
            outdir=outdir,
            with_sorting=True,
        )

    print("Total timing: ", time.time() - t)
