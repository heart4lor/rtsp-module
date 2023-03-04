./rtsp-simple-server
ffmpeg -f avfoundation -i ":0" -acodec libmp3lame -ab 32k -ac 1 -f rtsp rtsp://hk1.sunyongfei.cn:8554/audio/syf
ffmpeg -f avfoundation -i ":0" -acodec libmp3lame -ab 32k -ac 1 -f rtsp rtsp://localhost:8554/audio/syf
ffmpeg -i rtsp://localhost:8554/audio -c copy output.mp3


ffmpeg -f avfoundation -video_size 1280x720 -framerate 30 -i "1:0" -vcodec libx264 -preset ultrafast -acodec libmp3lame -f rtsp rtsp://0.0.0.0:8554/video
ffplay rtsp://hk1.sunyongfei.cn:8554/audio/syf