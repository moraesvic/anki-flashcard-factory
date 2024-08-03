#!/usr/bin/env bash

set -ex

cd src
go build .
cp ./flashcard-factory ~/bin/flashcard-factory

cat <<EOF > ~/bin/move-anki-audios
#!/usr/bin/env bash

set -ex

FLASHCARD_FOLDER="\$HOME/.local/share/Anki2/User 1/collection.media/"

move-anki-audios() {
    mv *.mp3 "\$FLASHCARD_FOLDER"
}

move-anki-audios
EOF

cat <<EOF > ~/bin/process-flashcards
#!/usr/bin/env bash

set -ex

# Replace this with your own AWS_PROFILE
export AWS_PROFILE="moraesvic"
cd \$(mktemp -d)

FLASHCARD_FOLDER="\$HOME/.local/share/Anki2/User 1/collection.media/"

flashcard-factory "\$HOME/anki_input.txt" > "\$HOME/anki_output.txt"

ls -la
move-anki-audios

cd -
EOF

chmod +x ~/bin/move-anki-audios
chmod +x ~/bin/process-flashcards

