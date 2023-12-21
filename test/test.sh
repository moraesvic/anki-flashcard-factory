# @TODO move this to the GO code

# When using "neural", the sample rate will already be the
# best one possible: https://awscli.amazonaws.com/v2/documentation/api/latest/reference/polly/synthesize-speech.html

aws polly synthesize-speech \
  --engine neural \
  --output-format mp3 \
  --voice-id Zhiyu --text "$(cat input.txt)" \
  polly.mp3

