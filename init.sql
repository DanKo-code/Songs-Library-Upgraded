CREATE DATABASE "songlibrarydb";
CREATE USER "songlibraryadmin" WITH PASSWORD 'SongLibraryAdmin2024';
GRANT ALL PRIVILEGES ON DATABASE "songlibrarydb" TO "songlibraryadmin";
\c "songlibrarydb";
GRANT ALL PRIVILEGES ON SCHEMA public TO "songlibraryadmin";