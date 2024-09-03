import os
import json
import random

import matplotlib.pyplot as plt
from matplotlib.colors import XKCD_COLORS as plot_colors

safelock_cmd = "safelock-cli"
pwd = "123456789"
rest = "60s"
input_path = "test"
output_name = "test"
output_dir = "safelock_dump"
runs = 3
figure_width = 14
figure_height = 2.5
bar_width = 0.6
measure = "Seconds"
root = os.getcwd()

def encrypt():
    err = os.system(
        f"hyperfine --runs {runs} --prepare "
        f"'sleep {rest}' "
        f"'echo \"{pwd}\" | {safelock_cmd} encrypt {input_path} {output_name}.sla --quiet' "
        f"'echo \"{pwd}\" | {safelock_cmd} encrypt {input_path} {output_name}_sha256.sla --quiet --sha256' "
        f"'echo \"{pwd}\" | {safelock_cmd} encrypt {input_path} {output_name}_sha512.sla --quiet --sha512' "
        f"'gpgtar -e -o test.gpg -c --yes --batch --gpg-args \"--passphrase {pwd}\" Videos/' "
        f"--export-json {root}/encryption.json"
    )

    if err:
        exit(err)

def decrypt():
    err = os.system(
        f"hyperfine --runs {runs} --prepare "
        f"'rm -rf {output_dir} {output_name}_1_ && mkdir {output_dir} && sleep {rest}' "
        f"'echo \"{pwd}\" | {safelock_cmd} decrypt {output_name}.sla {output_dir} --quiet' "
        f"'echo \"{pwd}\" | {safelock_cmd} decrypt {output_name}_sha256.sla {output_dir} --quiet --sha256' "
        f"'echo \"{pwd}\" | {safelock_cmd} decrypt {output_name}_sha512.sla {output_dir} --quiet --sha512' "
        f"'gpgtar -d --yes --batch --gpg-args \"--passphrase {pwd}\" test.gpg' "
        f"--export-json {root}/decryption.json"
    )

    if err:
        exit(err)

def get_label(i, clean=False):
    label = i['command']

    if 'gpg' in label:
        label = 'gpgtar'
    elif 'sha256' in label:
        label = 'safelock --sha256'
    elif 'sha512' in label:
        label = 'safelock --sha512'
    else:
        label = 'safelock'

    if clean:
        return label

    return f"{label}\n{i['median']:.3f}s"

# os.chdir(os.path.expanduser("~"))
# encrypt()
# decrypt()
os.chdir(root)

with open("encryption.json") as f:
    data = sorted(json.load(f)['results'], key=lambda i: i['median'])
    labels = [get_label(i) for i in data]
    scores = [i['median'] for i in data]
    colors_map = {get_label(i, 1): random.choice(list(plot_colors.values())) for i in data}
    colors = [colors_map[get_label(i, 1)] for i in data]

plt.margins(3.5)

fig, ax = plt.subplots()
ax.set_title('Encryption Time')
ax.set_xlabel(measure)
ax.barh(labels, scores, bar_width, color=colors)
fig.set_size_inches(w=figure_width, h=figure_height)
fig.tight_layout()
fig.savefig("encryption-time.webp", transparent=True, format="webp")

with open("decryption.json") as f:
    data = sorted(json.load(f)['results'], key=lambda i: i['median'])
    labels = [get_label(i) for i in data]
    decryption = [i['median'] for i in data]
    colors = [colors_map[get_label(i, 1)] for i in data]

fig, ax = plt.subplots()
ax.set_title('Decryption Time')
ax.set_xlabel(measure)
ax.barh(labels, decryption, bar_width, color=colors)
fig.set_size_inches(w=figure_width, h=figure_height)
fig.tight_layout()
fig.savefig("decryption-time.webp", transparent=True, format="webp")
