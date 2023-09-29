CreateTask
  TaskInit
    Force Pull Image or No Image
       Yes: PullImage
        ContainerCreate #create
       No: (#create)
        ContainerStartAndWait
          AutoRemove
            Yes: Delete stopped container
              (#done)
            No: Task Done #done
