docker build -t resume_ats .

docker rm -f resume_ats

docker run -d \
--name resume_ats \
--network ats-network \
--env-file .env \
-p 8080:8080 \
resume_ats