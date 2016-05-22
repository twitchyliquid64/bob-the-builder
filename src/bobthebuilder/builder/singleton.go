package builder


//For our application, we only need one Builder{}. We keep track of that here

var builderInstance *Builder

func init() {
  builderInstance = New()
}


func GetInstance()*Builder{
  return builderInstance
}
