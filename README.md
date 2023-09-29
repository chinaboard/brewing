# Brewing

Brewing provides a simple container-based task scheduling framework. 

The current implementation offers a solution for video text conversion. 

The process is as follows:

- Download video files based on a custom worker ([yt-dlp](https://github.com/yt-dlp/yt-dlp)).
- Extract audio content (ffmpeg).
- Convert audio to text via the [Whisper service](https://github.com/chinaboard/whisperX-service).
- Optimize text structure through [OpenAi](https://platform.openai.com/).
- Push results via the [Bark](https://github.com/Finb/bark-server) messaging service.

#### Demo: Atomic Expert Explains "Oppenheimer" Bomb Scenes | WIRED
- [Demo - Article](https://htmlpreview.github.io/?https://github.com/chinaboard/brewing/blob/master/docs/demo.html)
- [Demo - Timeline](https://htmlpreview.github.io/?https://github.com/chinaboard/brewing/blob/master/docs/demo-timeline.html)

### Video to Text
![Video to Text](https://github.com/chinaboard/brewing/blob/master/docs/assets/v2t-task-flow.png?raw=true)

### Basic Task workflow
![Basic task](https://github.com/chinaboard/brewing/blob/master/docs/assets/basic-task-flow.png?raw=true)

### Task in docker workflow
![Docker task](https://github.com/chinaboard/brewing/blob/master/docs/assets/docekr-task-flow.png?raw=true)
