include ./env


ifndef version 
$(error Missing var "name")
endif

ifndef version 
$(error Missing var "version")
endif


start:
	tmux new-session -d -s bc
	tmux setenv -t bc HISTFILE /dev/null
	tmux split-window -t "bc:0"   -v -p 50
	tmux send-keys -t "bc:0.0" "go run main.go --config conf/conf.toml" Enter
	tmux send-keys -t "bc:0.1" "curl -X POST http://127.0.0.1:3000/notify --data 'channels=hello,world&text=test' -v"
	tmux set-option -g mouse on
	tmux attach -t bc
	tmux kill-session -t bc


start-test:
	tmux new-session -d -s bc
	tmux setenv -t bc HISTFILE /dev/null
	tmux split-window -t "bc:0"   -v -p 50
	tmux send-keys -t "bc:0.0" "go run main.go --config conf/conf.toml.test" Enter
	tmux send-keys -t "bc:0.1" "curl -X POST http://127.0.0.1:3000/notify --data 'channels=hello,world&text=test' -v"
	tmux set-option -g mouse on
	tmux attach -t bc
	tmux kill-session -t bc




image:
	docker build -t $(name):latest . --build-arg goproxy=https://goproxy.cn,direct
	docker tag $(name):latest $(name):v$(version)
