go build -o bed-and-breakfast.exe ./cmd/web/. || exit /b
bed-and-breakfast.exe -dbname=bedandbreakfast -dbuser=postgres -dbpass=usman123 -cache=false production=false