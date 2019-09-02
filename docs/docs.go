// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-09-02 14:09:35.9496337 +0800 CST m=+0.135639501

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "HaoZhi.Cui",
            "url": "http://github.com/sak0",
            "email": "61755280@qq.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/chart": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取Chart仓库列表",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "PageSize",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "VerifyStatus",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Cached",
                        "name": "cached",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "chart_name",
                        "name": "chart",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "ClusterName",
                        "name": "cluster",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ChartRepo"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/chart/{chartName}/versions": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取指定Chart仓库的版本列表",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "PageSize",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "VerifyStatus",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Cached",
                        "name": "cached",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "ClusterName",
                        "name": "cluster",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ChartVersion"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/chart/{repo}/{version}": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "删除本地指定Chart仓库的指定版本",
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/models.ChartVersion"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/chart/{repo}/{version}/download": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "下载更新指定Chart仓库的指定版本到本地仓库",
                "parameters": [
                    {
                        "description": "Download the version to local",
                        "name": "version",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ChartVersion"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/models.ChartVersion"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/chart/{repo}/{version}/push": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "推送指定Chart仓库的指定版本到远端仓库",
                "parameters": [
                    {
                        "type": "string",
                        "description": "remote",
                        "name": "remote",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/models.ChartVersion"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/cluster": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取edge-cloud整体集群信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "PageSize",
                        "name": "pageSize",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SeederNode"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/health": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "获取服务健康状态",
                "responses": {
                    "200": {
                        "description": "{\"code\":S200,\"data\":{},\"msg\":\"ok\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/repository": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取镜像仓库列表",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "PageSize",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "VerifyStatus",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Cached",
                        "name": "cached",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Repository"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/repository/{repoName}/tags": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取单个镜像仓库的tag列表",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "PageSize",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "VerifyStatus",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Cached",
                        "name": "cached",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.RepositoryTag"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/repository/{repo}/{tag}": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "删除本地指定镜像仓库的指定tag",
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/models.RepositoryTag"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/repository/{repo}/{tag}/download": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "下载更新指定镜像仓库的指定tag到本地仓库",
                "parameters": [
                    {
                        "description": "Download the tag to local",
                        "name": "tag",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.RepositoryTag"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/models.RepositoryTag"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/versiondetail/file": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "查询指定Version的文件详情，例如：README.md",
                "parameters": [
                    {
                        "type": "string",
                        "description": "chart_name",
                        "name": "chart",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "version",
                        "name": "version",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "file_name",
                        "name": "file_name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ClusterName",
                        "name": "cluster",
                        "in": "query"
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/models.ChartVersion"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/versiondetail/filelist": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "查询指定Version的文件列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "chart_name",
                        "name": "chart",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "version",
                        "name": "version",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ClusterName",
                        "name": "cluster",
                        "in": "query"
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/models.ChartVersion"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/versiondetail/params": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "查询指定Version的参数Key-Value详情",
                "parameters": [
                    {
                        "type": "string",
                        "description": "chart_name",
                        "name": "chart",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "version",
                        "name": "version",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ClusterName",
                        "name": "cluster",
                        "in": "query"
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/models.ChartVersion"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ChartRepo": {
            "type": "object",
            "properties": {
                "cached": {
                    "type": "boolean"
                },
                "created": {
                    "type": "string"
                },
                "home": {
                    "type": "string"
                },
                "icon": {
                    "type": "string"
                },
                "latest_version": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner_name": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "updated": {
                    "type": "string"
                },
                "verify_status": {
                    "type": "string"
                },
                "version_count": {
                    "type": "integer"
                }
            }
        },
        "models.ChartVersion": {
            "type": "object",
            "properties": {
                "app_version": {
                    "type": "string"
                },
                "cached": {
                    "type": "boolean"
                },
                "created": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "digest": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "updated": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "verify_status": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "models.Repository": {
            "type": "object",
            "properties": {
                "cached": {
                    "type": "boolean"
                },
                "description": {
                    "type": "string"
                },
                "owner_name": {
                    "type": "string"
                },
                "pull_count": {
                    "type": "integer"
                },
                "repo_name": {
                    "type": "string"
                },
                "star_count": {
                    "type": "integer"
                },
                "tag_count": {
                    "type": "integer"
                },
                "verify_status": {
                    "type": "string"
                }
            }
        },
        "models.RepositoryTag": {
            "type": "object",
            "properties": {
                "architecture": {
                    "type": "string"
                },
                "author": {
                    "type": "string"
                },
                "cached": {
                    "type": "boolean"
                },
                "digest": {
                    "type": "string"
                },
                "docker_version": {
                    "type": "string"
                },
                "os": {
                    "type": "string"
                },
                "os_version": {
                    "type": "string"
                },
                "size": {
                    "type": "integer"
                },
                "tag_name": {
                    "type": "string"
                },
                "verify_status": {
                    "type": "string"
                }
            }
        },
        "models.SeederNode": {
            "type": "object",
            "properties": {
                "advertise_addr": {
                    "type": "string"
                },
                "bind_addr": {
                    "type": "string"
                },
                "chart_count": {
                    "type": "integer"
                },
                "cluster_name": {
                    "type": "string"
                },
                "image_count": {
                    "type": "integer"
                },
                "pull_count": {
                    "type": "integer"
                },
                "role": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1",
	Host:        "172.16.24.200:15000",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Seeder API",
	Description: "Server for image/chart repo consistent.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
