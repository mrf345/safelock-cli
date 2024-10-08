import os
import json
import random

import matplotlib.pyplot as plt
from matplotlib.colors import XKCD_COLORS as plot_colors

safelock_cmd = "~/Projects/safelock-cli/safelock-cli"
pwd = "123456789"
rest = "60s"
input_path = "Videos"
output_name = "test"
output_dir = "safelock_dump"
runs = 3
figure_width = 14
figure_height = 2.3
bar_width = 0.65
measure = "Seconds"
root = os.getcwd()

def get_label(i, key="command"):
    matchers = [
        ('gpg', 'gpgtar',),
        ('7z', '7zip (fastest)',),
        ('age', 'age (tar-zstd)'),
        ('safelock', 'safelock',),
    ]

    return next((v for m, v in matchers if m in i[key]))

def get_name(i):
    matchers = [
        ('gpg', f'{output_name}.gpg',),
        ('7z', f'{output_name}.7z',),
        ('age', f'{output_name}.age'),
        ('safelock', f'{output_name}.sla',),
    ]

    return next((v for m, v in matchers if m in i))

def encrypt():
    err = os.system(
        f"hyperfine --runs {runs} --prepare "
        f"'sleep {rest}' "
        f"'tar cv --zstd {input_path} | . {root}/pipe_age_password.sh | age -e -p -o {get_name('age')}' "
        f"'echo \"{pwd}\" | {safelock_cmd} encrypt {input_path} {get_name('safelock')} --quiet' "
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
        f"'sleep 0.05; xdotool type \"{pwd}\"; xdotool key \"Return\" | age --decrypt {get_name('age')} | tar x --zstd -f - -C {output_dir}' "
        f"'echo \"{pwd}\" | {safelock_cmd} decrypt {get_name('safelock')} {output_dir} --quiet' "
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
    colors_map = {get_label(i): random.choice(list(plot_colors.values())) for i in data}
    colors = [colors_map[get_label(i)] for i in data]

fig, ax = plt.subplots()
ax.set_title('Encryption Time')
ax.set_xlabel(measure)
ax.grid(zorder=0, axis='x')
ax.barh(labels, scores, bar_width, color=colors, zorder=3)
ax.bar_label(ax.containers[0], label_type='edge', padding=3, fmt=lambda i: f"{i:.2f}")
fig.set_size_inches(w=figure_width, h=figure_height)
fig.tight_layout()
fig.savefig("encryption-time.webp", transparent=True, format="webp")


# Decryption Time Plot

with open("decryption.json") as f:
    data = sorted(json.load(f)['results'], key=lambda i: i['median'])
    labels = [get_label(i) for i in data]
    decryption = [i['median'] for i in data]
    colors = [colors_map[get_label(i)] for i in data]

fig, ax = plt.subplots()
ax.set_title('Decryption Time')
ax.set_xlabel(measure)
ax.grid(zorder=0, axis='x')
ax.barh(labels, decryption, bar_width, color=colors, zorder=3)
ax.bar_label(ax.containers[0], label_type='edge', padding=3, fmt=lambda i: f"{i:.2f}")
fig.set_size_inches(w=figure_width, h=figure_height)
fig.tight_layout()
fig.savefig("decryption-time.webp", transparent=True, format="webp")


# File Sizes Plot

os.chdir(os.path.expanduser("~"))
data = sorted([{
    'size': os.path.getsize(get_name(get_label(i))) / 1024 / 1024,
    'label': get_label(i),
    'color': colors_map[get_label(i)],
} for i in data], key=lambda i: i['size'])
os.chdir(root)
labels = [get_label(i, 'label') for i in data]
sizes = [i['size'] for i in data]
colors = [i['color'] for i in data]

fig, ax = plt.subplots()
ax.set_title('File Size')
ax.set_xlabel("MegaBytes")
ax.grid(zorder=0, axis='x')
ax.barh(labels, sizes, bar_width, color=colors, zorder=3)
ax.bar_label(ax.containers[0], label_type='edge', padding=3, fmt=lambda i: f"{i:.0f}")
fig.set_size_inches(w=figure_width, h=figure_height)
fig.tight_layout()
fig.savefig("file-size.webp", transparent=True, format="webp")
