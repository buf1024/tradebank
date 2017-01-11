package ioms


// Connection wrap of net.Conn
type Connection interface {
    Listen()(error)
    Connect()(error)

}

// EndPoint tunnel endpoint
type EndPoint struct {

}

func ff()  {
    var con Connection;

    con.Connect
}