cp /home/blissnd/.config/sublime-text-3/Local/Session.sublime_session_template /home/blissnd/.config/sublime-text-3/Local/Session.sublime_session

find ./webserver -name "*.go" | xargs subl -n
find ./webserver -name "*.html" | xargs subl -n
find ./webserver -name "*.js" | xargs subl -n

