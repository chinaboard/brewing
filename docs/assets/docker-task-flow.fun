Task Init
  Image Not found      Force Pull
    Yes: PullImage
      ContainerCreate #create
    No: (#create)
      ContainerStartAndWait
        AutoRemove
          Yes: Delete stopped container
            (#done)
          No: Task Done #done
