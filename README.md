# VideoVerse API Documentation

Welcome to the documentation of the VideoVerse API. This API is developed using the [Go](https://go.dev/) programming language with the [Gin](https://gin-gonic.com/) web framework and leverages [GORM](https://gorm.io/) for ORM functionalities. It uses [PostgreSQL](https://www.postgresql.org/) as its primary database, ensuring fast, reliable, and scalable data handling. Designed to serve as a backend for a video-sharing platform similar to YouTube, this API supports video uploading, user engagement, playlist and channel management, and more.

## Features

- **Secure Authentication** is at the heart of the VideoVerse platform. Using [JWT](https://jwt.io/) for token-based auth and [bcrypt](https://en.wikipedia.org/wiki/Bcrypt) for password hashing, the system ensures user credentials are safely managed. Email verification is required for all accounts, and password reset flows are protected by time-limited tokens.

- **User Account Management** is designed to be complete and secure. Users can register, verify their email, log in, view and update their profile, and delete their accounts. Each profile supports rich user data including name, gender, birth date, description, profile picture, location, and linked social media accounts.

- **Video Uploading and Management** is streamlined through endpoints that support both image and video file uploads. Users can organize content into categories and playlists with customizable visibility settings (public, unlisted, private). Video categories and playlists can be updated or deleted, and videos include metadata like title, description, and status flags.

- **User Interactions** include the ability to like or dislike videos, comment with nested replies, and subscribe to channels. These actions are reflected in the notification system to keep users engaged and informed. Each reaction can be created or removed with authenticated calls.

- **Notification System** allows users to be alerted when key actions take place, such as likes, comments, or subscriptions. These notifications include redirect links to the appropriate resource for seamless UX on the frontend.

- **Commenting System** supports top-level comments as well as threaded replies. Each comment is associated with a video and a user, and supports display of the user profile, video context, and parent-child comment chains.

- **Filtering and Searching** is available across all GET endpoints. Results can be filtered by fields such as user ID, video ID, category, or visibility. Searching by keyword is also supported, and results can be sorted and ordered as needed to support frontend requirements like infinite scroll or pagination.

- **Robust Input Validation** ensures that every POST, PUT, or PATCH request is safe and consistent. Validation rules include required fields, proper UUID formats, valid emails, numeric constraints, string length limits, and boolean or enum enforcement where applicable.

- **API Key Security** is implemented via single-use HMAC-based API keys. This provides an extra layer of protection for every request beyond bearer tokens, helping to prevent abuse and unauthorized access. If needed, this security layer can be configured via environment variables.

- **SQL Injection Protection** is guaranteed by the use of parameterized queries in GORM, ensuring no user input is directly injected into SQL statements.

- **Comprehensive Documentation:** Clear [Postman](https://www.postman.com/) documentation for simplified testing.

## Additional Notes

### For Frontend Developers
- You can use `SPECIAL_API_KEY (Uh/UB%SKft3CU3e0zJAvBhp3cyo/un2021/zLQf1BKGZZuQ6w5P9VAM6Sj0CcQCm)`, put it directly in the http request header as `X-Api-Key`.
- Alternatively, if you want to try the One-Time API Key feature, the way to create the `X-Api-Key` are:
  1. Generate a random string.
  2. Calculate the HMAC signature between the random string and the `HMAC_KEY (dI62Fk_8wb2uL8CLmSLFkDoAO/tfDeod)` using SHA-256.
  3. The result of the HMAC calculation is combined with the random string with the pattern `random_string:hmac_result`.
  4. Then, encode the pattern with the Base64 algorithm, the endcode result is the `X-Api-key`.

### For Backend Developers
- if you don't want to use the  One-Time API key feature, don't forget to set `SPECIAL_API_KEY` in .env to the request header as `X-Api-Key` or change `ENABLE_API_KEY` in .env to `false` or you will not be able to access all endpoints at all.
- And don't forget to:
  - Crete [Cloudinary](https://cloudinary.com/) account for image uploads.
  - Create [Backblaze](https://www.backblaze.com/) account for file system needs.
  - Set up [Gmail SMTP](https://www.digitalocean.com/community/tutorials/how-to-use-google-s-smtp-server) for email sending.

---

© 2025 VideoVerse API Project. All rights reserved. By [fauzancodes](https://fauzancodes.id/)