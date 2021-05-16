import glob
from icu import Collator, Locale
import multiprocessing as mp
import os
import re
import shutil
import time
import yaml


def worker(path, outdir, with_sorting=True):
    collator = Collator.createInstance(Locale("pl_PL.UTF-8"))
    separator = re.compile("[\W\d]+")
    filepath = path.replace(".yml", ".txt")
    with open(filepath) as file:
        text = file.read().lower().rstrip()
        words = set(re.split(separator, text))
    with open(path) as file:
        meta = yaml.safe_load(file)
    with open(f"{outdir}/extracted-words-for-{meta['label']}.txt", "w") as file:
        if with_sorting:
            words = sorted(words, key=collator.getSortKey)
        file.write("\n".join(words))
    return path


if __name__ == "__main__":
    t = time.time()

    outdir = "words"
    if os.path.exists(outdir):
        shutil.rmtree(outdir)
    os.makedirs(outdir)

    pool = mp.Pool(mp.cpu_count())

    print("Processing")

    jobs = []
    for path in glob.glob("../data/pl/**/*.yml", recursive=True):
        process = mp.Process(target=worker, args=(path, outdir))
        jobs.append(process)
    for job in jobs:
        job.start()
    for job in jobs:
        job.join()

    print("Total timing: ", time.time() - t)
