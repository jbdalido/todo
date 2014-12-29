# TODO

**Has been made a million time, this is a personnal representation of one of the most made type of piece of software** 

Just another todo-list, in cli only, writing some big fat files in fat folders.
Use GDrive, DropBox, BitSync or whatever to sync across your devices.

The watch system could be simpler, it could use etcd for example as a backend storage, but i want it to be as portable as possible, PR are welcome.


# Init

Setup a folders as your todos repository :

    todo init /home/jbaptiste/DropBox/todos

Create a todo category

    todo create work

Setup a default category 

    todo default work

# Let's get to work

Add a task 

    todo add -m "a new todo" -c "work"

or

    echo "a new task" | todo add


Done a task

    todo done ae4Fr5
    ae4Fr5 is done ! 

Delete a task

    todo delete ae4Fr5
    ae4Fr5 is deleted, you're weak.

### TODO 

Watch a todo, and get slack notified :

    export SLACK_TOKEN=xcxc5234-5234-2345
    export SLACK_CHAN=general
    export SLACK_USERNAME=todo

    todo watch /home/jbaptiste/todos/ --hook="slack" --hook="os" -c "work" -c "private"

or with docker :

    docker run -v /home/jbaptiste/todos:/todos jbaptiste/todo watch --hook="slack" --hook="os" -c "work" -c "private"

Example Message:

    Todo
    jbaptiste has deleted task #aeFR5
    [work] "add a new task"


### END TODO

The task is saved as a json file looking like this :

    cat /home/jbaptiste/todos/work/ae4FR5

    {
        "id": "ae4Fr5",
        "user": "jbaptiste",
        "message": "a new task",
        "category": "work",
        "creationDate": 123413463456,
        "doneDate": 1234132424,
    }


# Command line

Here's a bunch a possibilities of what you can do with your todo. 

    todo -lo (ordered)
    - 15March
        | 1aef4f - [private] GoMarathon needs attention
        | mFQ4f - [work] Lorem
    - 12March
        | 1 - [private] - Loreum Dolor mothfucker

    todo -l
    # TODOS (3) :
        | 1 - GoMarathon needs attention
        | 2 - Lorem
        | 3 - Loreum Dolor mothfucker

    todo -lA (all) 
    # TODOS (3) :
        | 1 - GoMarathon needs attention
        | 2 - Lorem
        | 3 - Loreum Dolor mothfucker
    # DONE (343) (last 10) :
        | 1 - GoMarathon needs attention
        | 2 - Lorem
        | 3 - Loreum Dolor mothfucker

    todo -ld 15 (dones)
    # DONE (343) (last 15) :
        | 1 - GoMarathon needs attention
        | 2 - Lorem
        | 3 - Loreum Dolor mothfucker
        ....


    todo -L (Left)
    3

    todo -D (Done)
    0

    todo add -m "add a todo task" -c Work

# What's next

- Add encryption support
