import os
import json
import random

import matplotlib.pyplot as plt
from matplotlib.colors import XKCD_COLORS as plot_colors

safelock_cmd = "safelock-cli"
pwd = "123456789"
rest = "60s"
input_path = "~/Videos"
output_name = "test"
output_dir = "safelock_dump"
runs = 3
figure_width = 14
figure_height = 3
bar_width = 0.6
measure = "Seconds"
root = os.getcwd()

def get_label(i, clean=False, key="command"):
    matchers = [
        ('gpg', 'gpgtar',),
        ('7z', '7zip (fastest)',),
        ('256', 'safelock --sha256',),
        ('512', 'safelock --sha512',),
        ('safelock', 'safelock',),
    ]
    label = next((v for m, v in matchers if m in i[key]))

    if clean:
        return label
    if key == "label":
        return f"{label}\n{i['size']:.2f} MB"

    return f"{label}\n{i['median']:.3f}s"

def get_name(i):
    matchers = [
        ('gpg', f'{output_name}.gpg',),
        ('7z', f'{output_name}.7z',),
        ('256', f'{output_name}_sha256.sla',),
        ('512', f'{output_name}_sha512.sla',),
        ('safelock', f'{output_name}.sla',),
    ]

    return next((v for m, v in matchers if m in i))

def encrypt():
    err = os.system(
        f"hyperfine --runs {runs} --prepare "
        f"'sleep {rest}' "
        f"'echo \"{pwd}\" | {safelock_cmd} encrypt {input_path} {get_name('safelock')} --quiet' "
        f"'echo \"{pwd}\" | {safelock_cmd} encrypt {input_path} {get_name('256')} --quiet --sha256' "
        f"'echo \"{pwd}\" | {safelock_cmd} encrypt {input_path} {get_name('512')} --quiet --sha512' "
        f"'7z a -p{pwd} -mx1 {get_name('7z')} {input_path}' "
        f"'gpgtar -e -o {get_name('gpg')} -c --yes --batch --gpg-args \"--passphrase {pwd}\" {input_path}' "
        f"--export-json {root}/encryption.json"
    )

    if err:
        exit(err)

def decrypt():
    err = os.system(
        f"hyperfine --runs {runs} --prepare "
        f"'rm -rf {output_dir} {output_name}_*_ && mkdir {output_dir} && sleep {rest}' "
        f"'echo \"{pwd}\" | {safelock_cmd} decrypt {get_name('safelock')} {output_dir} --quiet' "
        f"'echo \"{pwd}\" | {safelock_cmd} decrypt {get_name('256')} {output_dir} --quiet --sha256' "
        f"'echo \"{pwd}\" | {safelock_cmd} decrypt {get_name('512')} {output_dir} --quiet --sha512' "
        f"'7z e -y -p{pwd} -mx1 {get_name('7z')} -o{output_dir}' "
        f"'gpgtar -d --yes --batch --gpg-args \"--passphrase {pwd}\" {get_name('gpg')}' "
        f"--export-json {root}/decryption.json"
    )

    if err:
        exit(err)

os.chdir(os.path.expanduser("~"))
encrypt()
decrypt()
os.chdir(root)
plt.margins(3.5)


# Encryption Time Plot

with open("encryption.json") as f:
    data = sorted(json.load(f)['results'], key=lambda i: i['median'])
    labels = [get_label(i) for i in data]
    scores = [i['median'] for i in data]
    colors_map = {get_label(i, 1): random.choice(list(plot_colors.values())) for i in data}
    colors = [colors_map[get_label(i, 1)] for i in data]

fig, ax = plt.subplots()
ax.set_title('Encryption Time')
ax.set_xlabel(measure)
ax.barh(labels, scores, bar_width, color=colors)
fig.set_size_inches(w=figure_width, h=figure_height)
fig.tight_layout()
fig.savefig("encryption-time.webp", transparent=True, format="webp")


# Decryption Time Plot

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


# File Sizes Plot

os.chdir(os.path.expanduser("~"))
data = sorted([{
    'size': os.path.getsize(get_name(get_label(i))) / 1024 / 1024,
    'label': get_label(i),
    'color': colors_map[get_label(i, 1)],
} for i in data], key=lambda i: i['size'])
os.chdir(root)
labels = [get_label(i, key='label') for i in data]
sizes = [i['size'] for i in data]
colors = [i['color'] for i in data]

fig, ax = plt.subplots()
ax.set_title('File Size')
ax.set_xlabel("Megabytes")
ax.barh(labels, sizes, bar_width, color=colors)
fig.set_size_inches(w=figure_width, h=figure_height)
fig.tight_layout()
fig.savefig("file-size.webp", transparent=True, format="webp")
