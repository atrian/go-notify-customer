// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/events": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Event"
                ],
                "summary": "Запрос всех доступных шаблонов",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.Event"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Event"
                ],
                "summary": "сохранение бизнес события",
                "parameters": [
                    {
                        "description": "Принимает dto события, отдает сохраненное событие с идентификатором",
                        "name": "event",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.IncomingEvent"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.Event"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/v1/events/{event_uuid}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Event"
                ],
                "summary": "Запрос деталей бизнес события",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID события в формате UUID v4",
                        "name": "event_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.Event"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Event"
                ],
                "summary": "обновление бизнес события",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID бизнес события в формате UUID v4",
                        "name": "event_uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Test",
                        "name": "event",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.IncomingEvent"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.Event"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Event"
                ],
                "summary": "удаление бизнес события",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID бизнес события в формате UUID v4",
                        "name": "event_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/v1/notifications": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Notifications"
                ],
                "summary": "отправка уведомлений",
                "parameters": [
                    {
                        "description": "Принимает JSON dto уведомлений, возвращает код 200 при успешной постановке, 429 при привышении лимита",
                        "name": "notification",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.IncomingNotification"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "429": {
                        "description": "Too Many Requests"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/v1/notifications/seed": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Notifications"
                ],
                "summary": "Создает бизнес событие и шаблон к нему, возвращает подготовленный JSON для запроса POST /api/v1/notifications",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.IncomingNotification"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/v1/stats": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Stat"
                ],
                "summary": "Запрос всей статистики",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.Stat"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/v1/stats/notification/{notification_uuid}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Stat"
                ],
                "summary": "Запрос статистики отправок по уведомлению",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID уведомления в формате UUID v4",
                        "name": "notification_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.Stat"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/v1/stats/person/{person_uuid}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Stat"
                ],
                "summary": "Запрос статистики отправок по пользователю",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID пользователя в формате UUID v4",
                        "name": "person_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.Stat"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/v1/templates": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Template"
                ],
                "summary": "Запрос деталей шаблона сообщения",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.Template"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Template"
                ],
                "summary": "сохранение шаблона сообщения",
                "parameters": [
                    {
                        "description": "Принимает dto нового шаблона сообщения, возвращает JSON сохраненными данными и идентификатором",
                        "name": "template",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.IncomingTemplate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.Template"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/v1/templates/{template_uuid}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Template"
                ],
                "summary": "Запрос деталей шаблона сообщения",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID шаблона в формате UUID v4",
                        "name": "template_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.Template"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Template"
                ],
                "summary": "обновление шаблона сообщения",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID шаблона сообщения в формате UUID v4",
                        "name": "template_uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Принимает dto шаблона сообщения, возвращает JSON с обновленными данными",
                        "name": "template",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.IncomingTemplate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.Template"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Template"
                ],
                "summary": "удаление шаблона сообщения",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID шаблона в формате UUID v4",
                        "name": "template_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.Event": {
            "type": "object",
            "properties": {
                "default_priority": {
                    "description": "DefaultPriority приоритет уведомления с таким событием по умолчанию",
                    "type": "integer"
                },
                "description": {
                    "description": "Description описание бизнес события",
                    "type": "string"
                },
                "event_uuid": {
                    "description": "EventUUID связь с UUID бизнес события",
                    "type": "string"
                },
                "notification_channels": {
                    "description": "NotificationChannels каналы отправки для данного события",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "description": "Title название бизнес события",
                    "type": "string"
                }
            }
        },
        "dto.IncomingEvent": {
            "type": "object",
            "properties": {
                "default_priority": {
                    "description": "DefaultPriority приоритет уведомления с таким событием по умолчанию",
                    "type": "integer"
                },
                "description": {
                    "description": "Description описание бизнес события",
                    "type": "string"
                },
                "notification_channels": {
                    "description": "NotificationChannels каналы отправки для данного события",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "description": "Title название бизнес события",
                    "type": "string"
                }
            }
        },
        "dto.IncomingNotification": {
            "type": "object",
            "properties": {
                "event_uuid": {
                    "description": "EventUUID связь с UUID бизнес события",
                    "type": "string"
                },
                "message_params": {
                    "description": "MessageParams key-value подстановки в шаблон уведомления",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.MessageParam"
                    }
                },
                "person_uuids": {
                    "description": "PersonUUIDs связь с пользователями - получателями уведомления",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "priority": {
                    "description": "Priority опциональный приоритет уведомления",
                    "type": "integer"
                }
            }
        },
        "dto.IncomingTemplate": {
            "type": "object",
            "properties": {
                "body": {
                    "description": "Body тело шаблона",
                    "type": "string"
                },
                "channel_type": {
                    "description": "ChannelType связь с каналом отправки",
                    "type": "string"
                },
                "description": {
                    "description": "Description описание шаблона",
                    "type": "string"
                },
                "event_uuid": {
                    "description": "EventUUID связь с UUID бизнес события",
                    "type": "string"
                },
                "title": {
                    "description": "Title название шаблона",
                    "type": "string"
                }
            }
        },
        "dto.MessageParam": {
            "type": "object",
            "properties": {
                "key": {
                    "description": "Key ключ по которому будет произведен поиск в теле уведомления",
                    "type": "string"
                },
                "value": {
                    "description": "Value значение которое будет подставлено вместо ключа в шаблоне",
                    "type": "string"
                }
            }
        },
        "dto.Stat": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "CreatedAt дата и время отправки",
                    "type": "string"
                },
                "notification_uuid": {
                    "description": "NotificationUUID связь с уведомлением",
                    "type": "string"
                },
                "person_uuid": {
                    "description": "PersonUUID связь отправленного уведомления с клиентом",
                    "type": "string"
                },
                "stat_uuid": {
                    "description": "StatUUID id записи статистики",
                    "type": "string"
                },
                "status": {
                    "description": "Status статус отправки",
                    "allOf": [
                        {
                            "$ref": "#/definitions/dto.StatStatus"
                        }
                    ]
                }
            }
        },
        "dto.StatStatus": {
            "type": "integer",
            "enum": [
                1,
                2,
                3
            ],
            "x-enum-comments": {
                "BadChannel": "Канал отправки не поддерживается",
                "Failed": "Ошибка отправки",
                "Sent": "Уведомление отправлено"
            },
            "x-enum-varnames": [
                "Sent",
                "Failed",
                "BadChannel"
            ]
        },
        "dto.Template": {
            "type": "object",
            "properties": {
                "body": {
                    "description": "Body тело шаблона",
                    "type": "string"
                },
                "channel_type": {
                    "description": "ChannelType связь с каналом отправки",
                    "type": "string"
                },
                "description": {
                    "description": "Description описание шаблона",
                    "type": "string"
                },
                "event_uuid": {
                    "description": "EventUUID связь с UUID бизнес события",
                    "type": "string"
                },
                "template_uuid": {
                    "description": "TemplateUUID - id шаблона",
                    "type": "string"
                },
                "title": {
                    "description": "Title название шаблона",
                    "type": "string"
                }
            }
        }
    },
    "tags": [
        {
            "description": "\"Группа запросов бизнес событий\"",
            "name": "Event"
        },
        {
            "description": "\"Группа запросов уведомлений\"",
            "name": "Notifications"
        },
        {
            "description": "\"Группа запросов для работы со статистикой отправки\"",
            "name": "Stat"
        },
        {
            "description": "\"Группа запросов для работы с шаблонами сообщений. Для создания сообщений требуется предварительное создание бизнес событий\"",
            "name": "Template"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Go-notify-client",
	Description:      "Сервис отправки уведомлений.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
