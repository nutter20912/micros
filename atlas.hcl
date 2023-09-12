variable "service" {
  type        = string
  description = "service 名稱"
}

data "external_schema" "gorm" {
  program = [
    "atlas-provider-gorm",
    "load",
    "--path", "./app/${var.service}/models",
    "--dialect", "mysql",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://mysql/8/dev"
  migration {
    dir = "file://database/migrations/${var.service}"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
