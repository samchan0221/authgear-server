module github.com/authgear/authgear-server

go 1.13

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/Masterminds/squirrel v1.5.0
	github.com/authgear/graphql-go-relay v0.0.0-20201016065100-df672205b892
	github.com/boombuler/barcode v1.0.1
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/elastic/go-elasticsearch/v7 v7.13.1
	github.com/getsentry/sentry-go v0.11.0
	github.com/go-http-utils/etag v0.0.0-20161124023236-513ea8f21eb1
	github.com/go-http-utils/fresh v0.0.0-20161124030543-7231e26a4b27 // indirect
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a // indirect
	github.com/go-redis/redis/v8 v8.10.0
	github.com/go-redsync/redsync/v4 v4.3.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.2.0
	github.com/google/wire v0.5.0
	github.com/gorilla/csrf v1.7.0
	github.com/graphql-go/graphql v0.7.9
	github.com/graphql-go/handler v0.2.3
	github.com/iawaknahc/gomessageformat v0.0.0-20210428033148-c3f8592094b5
	github.com/iawaknahc/jsonschema v0.0.0-20201115095512-87990d0baba1
	github.com/iawaknahc/originmatcher v0.0.0-20200622040912-c5bfd3560192
	github.com/jetstack/cert-manager v1.4.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/joho/godotenv v1.3.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	// jwx >= 1.2.1 fix the bug that `alg` and `use` are NOT required.
	github.com/lestrrat-go/jwx v1.2.1
	github.com/lib/pq v1.10.2
	github.com/lithdew/quickjs v0.0.0-20200714182134-aaa42285c9d2
	github.com/njern/gonexmo v2.0.0+incompatible
	github.com/nyaruka/phonenumbers v1.0.70
	github.com/oschwald/geoip2-golang v1.5.0
	github.com/pquerna/otp v1.3.0
	// The changes are compatible. See https://github.com/rubenv/sql-migrate/compare/8d140a17f351..55d5740dbbcc
	github.com/rubenv/sql-migrate v0.0.0-20210614095031-55d5740dbbcc
	// The changes are compatible. See https://github.com/sfreiberg/gotwilio/compare/169c4cd5c691..c426a3710ab5
	github.com/sfreiberg/gotwilio v0.0.0-20201211181435-c426a3710ab5
	github.com/sirupsen/logrus v1.8.1
	github.com/skygeario/go-confusable-homoglyphs v0.0.0-20191212061114-e2b2a60df110
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.0
	github.com/test-go/testify v1.1.4 // indirect
	github.com/trustelem/zxcvbn v1.0.1
	// The changes are compatible. See https://github.com/ua-parser/uap-go/compare/e1c09f13e2fe..347a3497cc39
	github.com/ua-parser/uap-go v0.0.0-20210121150957-347a3497cc39
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/text v0.3.6
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/h2non/gock.v1 v1.1.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	nhooyr.io/websocket v1.8.7
	sigs.k8s.io/yaml v1.2.0
)
