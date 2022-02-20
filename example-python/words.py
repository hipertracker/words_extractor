import argparse
import glob
from typing import Tuple

try:
    from icu import Collator, Locale

    i18nsorting = True
except ModuleNotFoundError:
    # Not supported by M1
    i18nsorting = False

import multiprocessing as mp
import os
import re
import shutil
import time
import yaml


def worker(path: str, outdir: str, sorting: bool = False) -> Tuple[str, int]:
    # if sorting:
    #     if i18nsorting:
    #         collator = Collator.createInstance(Locale("pl_PL.UTF-8"))
    #     print("I18nN sorting not available")

    separator = re.compile("[\W\d]+")
    filepath = path.replace(".yml", ".txt")
    filesize = os.path.getsize(filepath)
    with open(filepath) as file:
        text = file.read().lower().rstrip()
        words = set(re.split(separator, text))
        try:
            words.remove('')
        except KeyError:
            pass
        words = list(words)
    with open(path) as file:
        meta = yaml.safe_load(file)
    with open(f"{outdir}/{meta['lang']}-{meta['code']}.txt", "w") as file:
        if sorting:
            if i18nsorting:
                words = sorted(words, key=collator.getSortKey)
            else:
                words.sort()
        file.write("\n".join(words))
    return path, filesize


if __name__ == "__main__":
    program_name = os.path.basename(__file__)
    cores = mp.cpu_count()
    parser = argparse.ArgumentParser(f'python {program_name}')
    parser.add_argument('-n', type=int, help=f'Number of cores to run (default: {cores})', default=cores)
    parser.add_argument('-s', '--sort', action='store_true', help='Sort results')
    args = parser.parse_args()
    if not 1 <= args.n <= cores:
        args.n = 10

    cpu_cores = args.n
    sorting = args.sort

    t = time.time()

    outdir = "words"
    src_path = "../data/??/**/*.yml"

    if os.path.exists(outdir):
        shutil.rmtree(outdir)
    os.makedirs(outdir)

    paths = glob.glob(src_path, recursive=True)
    if not paths:
        raise Exception(f"WRONG PATH {src_path}")

    print(f"Running using {cpu_cores} processes", end='')
    if sorting:
        if i18nsorting:
            print(" with sorting using collations")
        else:
            print(" with sorting")
    results: list = []
    pool = mp.Pool(cpu_cores)
    for path in paths:
        res = pool.apply_async(
            worker,
            kwds=dict(
                path=path,
                outdir=outdir,
                sorting=sorting,
            ),
        )
        results.append(res)
    total_size = 0
    items_count = len(results)
    for i, res in enumerate(results):
        path, size = res.get()
        total_size += size
        print(f"[{i + 1}/{items_count}] {path}")
    print(f"Total files: {items_count}")
    print(f"Total size: {round((total_size / 1024 / 1024))} MB")
    t = time.time() - t
    print(f"Total time: {t:.4f} s")
