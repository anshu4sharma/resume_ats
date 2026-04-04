docker build -t resume_ats . && \
docker run -d \
--name resume_ats \
--restart unless-stopped \
-p 8080:8080 \
--env-file .env \
resume_ats