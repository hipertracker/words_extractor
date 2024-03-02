import argparse
import glob
import multiprocessing as mp
import os
import shutil
import time
from typing import NoReturn

import regex
import yaml
from alive_progress import alive_bar
from icecream import ic
from icu import Collator, Locale

LOCALES: dict[str, str] = {
    "ar": "ar_SA.UTF-8",  # Arabic, Saudi Arabia
    "cs": "cs_CZ.UTF-8",  # Czech, Czech Republic
    "da": "da_DK.UTF-8",  # Danish, Denmark
    "de": "de_DE.UTF-8",  # German, Germany
    "el": "el_GR.UTF-8",  # Greek, Greece
    "en": "en_EN.UTF-8",  # English, not country-specific
    "eo": "eo.UTF-8",  # Esperanto, not country-specific
    "es": "es_ES.UTF-8",  # Spanish, Spain
    "fi": "fi_FI.UTF-8",  # Finnish, Finland
    "fr": "fr_FR.UTF-8",  # French, France
    "he": "he_IL.UTF-8",  # Hebrew, Israel
    "hr": "hr_HR.UTF-8",  # Croatian, Croatia
    "hu": "hu_HU.UTF-8",  # Hungarian, Hungary
    "it": "it_IT.UTF-8",  # Italian, Italy
    "lt": "lt_LT.UTF-8",  # Lithuanian, Lithuania
    "la": "en_EN.UTF-8",  # Latin use the same characters as English
    "nl": "nl_NL.UTF-8",  # Dutch, Netherlands
    "pl": "pl_PL.UTF-8",  # Polish, Poland
    "pt": "pt_PT.UTF-8",  # Portuguese, Portugal
    "ru": "ru_RU.UTF-8",  # Russian, Russia
    "sk": "sk_SK.UTF-8",  # Slovak, Slovakia
    "sv": "sv_SE.UTF-8",  # Swedish, Sweden
    "uk": "uk_UA.UTF-8",  # Ukrainian, Ukraine
}


def extract_words(path: str, lang: str, sorting: bool = False) -> list[str]:
    """Extract sorted (optional) unique utf-8 words from the file"""
    separator = regex.compile(r"\P{L}+")
    words: list[str] = []
    with open(path, "r", encoding="utf-8") as file:
        for line in file.readlines():
            _line = " ".join(line.strip().lower().split(" ")[2:-1])
            words += (w for w in set(regex.split(separator, _line)) if w and len(w) > 1)
    unique_words = list(set(words))
    if sorting:
        locale = LOCALES[lang]
        collator = Collator.createInstance(Locale(locale))
        return sorted(unique_words, key=collator.getSortKey)
    return unique_words


def worker(out_dir: str, yaml_filepath: str, sorting: bool) -> (str, int):
    with open(yaml_filepath, "r", encoding="utf-8") as file:
        meta = yaml.safe_load(file)

    lang = meta.get("lang")
    if lang is None:
        raise RuntimeError(f"Language not found in {yaml_filepath}")

    code = meta.get("code")
    if code is None:
        raise RuntimeError(f"Code not found in {yaml_filepath}")

    words = extract_words(
        path=yaml_filepath.replace("yml", "txt"),
        lang=lang,
        sorting=sorting,
    )

    out_filepath = f"{out_dir}/{lang}-{code}.txt"
    with open(out_filepath, "w", encoding="utf-8") as file:
        file.write("\n".join(words))

    return out_filepath, os.path.getsize(out_filepath)


def main(src: str, out_dir: str) -> NoReturn:
    program_name = os.path.basename(__file__)
    cores = mp.cpu_count()
    parser = argparse.ArgumentParser(f"python {program_name}")
    parser.add_argument(
        "-n", type=int, help=f"Number of cores to run (default: {cores})", default=cores
    )
    parser.add_argument("-s", "--sort", action="store_true", help="Sort results")
    args = parser.parse_args()
    if not 1 <= args.n <= cores:
        args.n = 10
    cpu_cores = args.n
    sorting = args.sort

    if os.path.exists(out_dir):
        shutil.rmtree(out_dir)
    os.makedirs(out_dir)

    print("Processing")

    print(f"Running using {cpu_cores} processes", end="")
    if sorting:
        print(" with sorting")


    t = time.time()

    paths = glob.glob(src, recursive=True)
    if not paths:
        raise RuntimeError(f"Wrong path pattern {src}")

    # spawn processes
    pool = mp.Pool(cpu_cores)
    results: list = []
    for yaml_filepath in paths:
        kwargs = {"yaml_filepath": yaml_filepath, "out_dir": out_dir, "sorting": sorting}
        res = pool.apply_async(worker, kwds=kwargs)
        results.append(res)

    # collect results
    items_count = len(results)
    total_size = 0
    with alive_bar(items_count) as bar:
        for res in results:
            total_size += res.get()[1]
            bar()

    print(f"Total files: {items_count}")
    print(f"Total size: {round((total_size / 1024 / 1024))} MB")
    t = time.time() - t
    print(f"Total time: {t:.4f} s")

if __name__ == "__main__":
    main(src="../data/??/**/*.yml", out_dir="words")
