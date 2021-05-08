import glob
from icu import Collator, Locale
import os
import re
import yaml
import shutil

with_sorting = False

outdir = "words"
if os.path.exists(outdir):
    shutil.rmtree(outdir)
os.makedirs(outdir)

collator = Collator.createInstance(Locale("pl_PL.UTF-8"))
separator = re.compile("[\W\d]+")
for path in glob.glob("../data/pl/**/*.yml", recursive=True):
    print(os.path.basename(path))
    with open(path) as file:
        meta = yaml.safe_load(file)
    with open(path.replace(".yml", ".txt")) as file:
        text = file.read().lower().rstrip()
        words = set(re.split(separator, text))
    with open(f"{outdir}/extracted-words-for-{meta['label']}.txt", "w") as file:
        if with_sorting:
            words = sorted(words, key=collator.getSortKey)
        file.write("\n".join(words))
