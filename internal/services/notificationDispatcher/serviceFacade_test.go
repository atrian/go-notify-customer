package notificationDispatcher

import (
	"context"
	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/google/uuid"
	"testing"
)

func TestServiceFacade_prepareTemplate(t *testing.T) {
	type fields struct {
		contact  contactVault
		template interface {
			FindByEventId(ctx context.Context, eventUUID uuid.UUID) ([]dto.Template, error)
		}
		event interface {
			FindById(ctx context.Context, eventUUID uuid.UUID) (dto.Event, error)
		}
	}
	type args struct {
		template string
		replaces []dto.MessageParam
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Regexp",
			args: args{
				template: "message [param1] and [param2] data",
				replaces: []dto.MessageParam{
					{
						Key:   "param1",
						Value: "VAL1",
					}, {
						Key:   "param2",
						Value: "VAL2",
					},
				},
			},
			want: "message VAL1 and VAL2 data",
		}, {
			name: "empty replaces. replace multi spaces with one",
			args: args{
				template: "message [param1] and [param2] data",
				replaces: []dto.MessageParam{},
			},
			want: "message and data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ServiceFacade{
				contact:  tt.fields.contact,
				template: tt.fields.template,
				event:    tt.fields.event,
			}
			if got := f.prepareTemplate(tt.args.template, tt.args.replaces); got != tt.want {
				t.Errorf("prepareTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
