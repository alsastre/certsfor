package client

// Status returns status information if connection is ok to the API, error otherwise
func (c *Client) Status() (status APIStatus, err error) {

	_, err = c.http.Get("/status").ReceiveSuccess(&status)
	return

}