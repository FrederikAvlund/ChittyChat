How to use Chitty Chat service


	-- If the server makes

	1)	Open a minimum of 4 different terminals at the folder Chitty-Chat.
	2)	1 of these terminals has to locate to the server folder
	3)	The rest has to locate to the client folder
	4)	From within the server folder, run go run server.go
	⁃	You should see the message: “Starting server at port :8080”
	5)	Now that the server is open, head to the remaining terminals and type in one of the following:
	⁃	go run client.go 
	⁃	This is if you want to try the service using an Anonymous name
	⁃	go run client.go -N Testbruger1
	⁃	Here, a participant called Testbruger1 is created.
	6)	Do step number 5 for each of the remaining terminals and once you have chosen a name, type Join in the terminal. Now, you have entered the Chitty-Chat and can communicate between the rest of the participants.
	7)	When you are ready to leave, type Leave in the terminal, and you will no longer receive any messages that are being broadcasted by the server.
