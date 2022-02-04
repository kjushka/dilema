package action

type FirstService interface {
    MethodOne(ctx dilema.Context, val int) (int, error)
    MethodTwo(struct {
        fs FirstService
        ss SecondService
    })
}

type SecondService interface {
    MethodThree() string
}