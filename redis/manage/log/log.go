package lib

type ServiceLogger interface {
	Logger
	Debug()
	Debugln()
	Debugf()
	EnableDebug()
	DisableDebug()
}

type Logger interface {
	Print(...interface{})
	Println(...interface{})
	Printf(string, ...interface{})
}

type RedisLog struct {
	debugOn bool
	Logger
}

func (rl *RedisLog) Debug(v ...interface{}) {
	if rl.debugOn {
		rl.Print(v)
	}
}

func (rl *RedisLog) Debugln(v ...interface{}) {
	if rl.debugOn {
		rl.Println(v)
	}
}

func (rl *RedisLog) Debugf(format string, v ...interface{}) {
	if rl.debugOn {
		rl.Printf(format, v)
	}
}

func (rl *RedisLog) EnableDebug() {
	rl.debugOn = true
	rl.Println("Debug logging enabled")
}

func (rl *RedisLog) DisableDebug() {
	rl.debugOn = true
	rl.Println("Debug logging disabled")
}
