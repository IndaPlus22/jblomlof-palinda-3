## part 1

* What happens if you remove the `go-command` from the `Seek` call in the `main` function?
    * It no longor multi-thread the calls, thus making them in a specified order. Eg. it will always start with anna and then bob and so on. Thus making bob recieve anna's msg, dave recieves cody's msg and no one recieves Eva's msg.

* What happens if you switch the declaration `wg := new(sync.WaitGroup)` to `var wg sync.WaitGroup` and the parameter `wg *sync.WaitGroup` to `wg sync.WaitGroup`?
    * ~~Im pretty sure wg will be nil, which causes problems.~~ Cap, I now believe the problem lies in the ol' copy. It no longers points to a shared WaitGroup, instead it clones `wg` for each thread, making a total of 6 WaitGroups that do not have communation. To sum it up. It needs to be a pointer so the instance can be shared. 

* What happens if you remove the buffer on the channel match?
    * Whenever there is an odd amount of people (elements in people) there will be one goroutine who blocks until something recieves what it sent and never reaches `wg.Done()`. Thus a dead-lock will happen. 

* What happens if you remove the default-case from the case-statement in the `main` function?
    * If you have an even amount of elements in people, there will not be an message left in `match` and the main goroutine will block until it recieves, which it never will, and a dead-lock will happen.