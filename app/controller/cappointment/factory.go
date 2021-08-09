package cappointment

import (
	"github.com/gin-gonic/gin"
)

type IFactory interface {
	IAppointment
}

type AppointmentFactory struct {
}

func (factory *AppointmentFactory) Create(appointmentType int, c *gin.Context) IFactory {
	if appointmentType < 0 {
		return nil
	}

	switch appointmentType {
	case 0:
		return NewVenue(c)
	case 1:
		return NewCoach(c)
	case 2:
		return NewCourse(c)
	}

	return nil
}
