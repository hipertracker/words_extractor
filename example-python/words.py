import glob
from typing import Tuple
import os

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
    if sorting:
        if i18nsorting:
            collator = Collator.createInstance(Locale("pl_PL.UTF-8"))
        print("I18nN sorting not available")

    separator = re.compile("[\W\d]+")
    filepath = path.replace(".yml", ".txt")
    filesize = os.path.getsize(filepath)
    with open(filepath) as file:
        text = file.read().lower().rstrip()
        words = set(re.split(separator, text))
    with open(path) as file:
        meta = yaml.safe_load(file)
    with open(f"{outdir}/{meta['lang']}-{meta['code']}.txt", "w") as file:
        if sorting and i18nsorting:
            words = sorted(words, key=collator.getSortKey)
        file.write("\n".join(words))
    return path, filesize


if __name__ == "__main__":
    t = time.time()

    outdir = "words"
    src_path = "../data/??/**/*.yml"

    if os.path.exists(outdir):
        shutil.rmtree(outdir)
    os.makedirs(outdir)

    pool = mp.Pool(mp.cpu_count())

    print("Processing")
    results = []
    paths = glob.glob(src_path, recursive=True)
    if not paths:
        raise Exception(f"WRONG PATH {src_path}")

    for path in paths:
        res = pool.apply_async(
            worker,
            kwds=dict(
                path=path,
                outdir=outdir,
                sorting=False,
            ),
        )
        results.append(res)
    total_size = 0
    items_count = len(results)
    for i, res in enumerate(results):
        path, size = res.get()
        total_size += size
        print(f"[{i+1}/{items_count}] {path}")
    print(f"Total files: {items_count}")
    print(f"Total size: {round((total_size / 1024 / 1024))} MB")
    t = time.time() - t
    print(f"Total time: {t:.4f} s")

