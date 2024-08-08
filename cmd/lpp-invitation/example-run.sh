
export CONTENT_DB_CONNECTION_STR="<db-url>"
export CONTENT_DB_USERNAME="<db-username>"
export CONTENT_DB_PASSWORD="<db-password>"
export CONTENT_DB_CONNECTION_PREFIX="<+srv or empty>"
export DB_TIMEOUT=30
export DB_IDLE_CONN_TIMEOUT=45
export DB_MAX_POOL_SIZE=8
export DB_DB_NAME_PREFIX="<db-name-prefix>"

export INSTANCE_IDS="tekenradar"

export CSV_PATH="participants.csv"
export ENV_SEPARATOR=";"
export FORCE_REPLACE="false"
export INVITATION_EMAIL_TEMPLATE_PATH="lpp-invitation.html"
export INVITATION_EMAIL_SUBJECT="Uitnodiging voor deelname aan Tekenradar"

export RUN_PARTICIPANT_CREATION="true"
export RUN_INVITATION_SENDING="true"

export EMAIL_CLIENT_URL="localhost:5005"

go run *.go