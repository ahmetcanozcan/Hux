package sergo

import (
	"testing"
)

func TestBasicEvent(t *testing.T){
		// Create socket
		s := newSocket(nil)
		done := make(chan bool)

		go func(){
			msg := <-s.GetEvent("foo")
			t.Log(msg)
			done<-true
		}()

		s.GetEvent("foo")<-"Hi"
		<-done
		
}

func TestEventsWithSelect(t *testing.T) {
	done := make(chan bool)
	s := newSocket(nil)

	go func(){
		for i := 0; i<3; {
			select{
			case msg := <- s.GetEvent("foo"):
				t.Log("foo:",msg)
				i++
			case msg := <- s.GetEvent("bar"):
				t.Log("bar:",msg)
				i++
			case msg := <- s.GetEvent("joe"):
				t.Log("joe:",msg)
				i++
			default:
				continue
			}
		}
		done<-true
	}()
	s.GetEvent("foo")<-"hi"
	s.GetEvent("bar")<-"hi"
	s.GetEvent("joe")<-"hi"
	

	<-done
}
