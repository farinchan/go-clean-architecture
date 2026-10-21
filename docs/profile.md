# Dokumentasi Lengkap API Profile - Nagastra Pay

Dokumen ini menyediakan panduan integrasi dan spesifikasi teknis lengkap untuk grup endpoint **Profile** pada backend Nagastra Pay.

---

## 1. Arsitektur & Aturan Autentikasi

Seluruh endpoint pada grup Profile berada di bawah proteksi middleware autentikasi JWT (`middleware.AuthMiddleware`):

- **Header Wajib:**
  ```http
  Authorization: Bearer <JWT_ACCESS_TOKEN>
  ```
  Token diperoleh dari endpoint login (`POST /api/v1/auth/login`) atau verifikasi 2FA (`POST /api/v1/auth/verify-2fa`).
- **Context User:**
  Middleware akan memvalidasi tanda tangan JWT, mengecek masa aktif (`exp`), dan mengekstrak `userID`, `email`, serta `role` ke dalam context request.
- **Notifikasi Keamanan Otomatis:**
  Pergantian password melalui endpoint `PATCH /api/v1/profile/change-password` akan memicu pengiriman email konfirmasi keamanan secara asynchronous ke alamat email terdaftar pengguna.
- **Format Respon Standar:**
  Seluruh endpoint mematuhi envelope standar `pkg/response`:
  ```json
  {
    "success": true,
    "message": "Pesan status operasi",
    "data": { ... },
    "error": null,
    "meta": null
  }
  ```

---

## 2. Ringkasan Endpoint Profile

| No | Method | Endpoint Path | Header Wajib | Deskripsi Singkat |
|:---|:---|:---|:---|:---|
| 1 | `GET` | `/api/v1/profile` | `Bearer <token>` | Mendapatkan data profil pengguna yang sedang login |
| 2 | `PUT` | `/api/v1/profile` | `Bearer <token>` | Memperbarui informasi profil (nama, email, no. HP, 2FA) |
| 3 | `PATCH` | `/api/v1/profile/2fa` | `Bearer <token>` | Mengaktifkan atau menonaktifkan fitur 2FA via email |
| 4 | `PATCH` | `/api/v1/profile/photo` | `Bearer <token>` | Mengunggah dan memperbarui foto profil pengguna |
| 5 | `PATCH` | `/api/v1/profile/change-password` | `Bearer <token>` | Mengganti password akun terautentikasi |
| 6 | `POST` | `/api/v1/profile/notification-secret/regenerate` | `Bearer <token>` | Membuat ulang secret notifikasi webhook (whsec_...) |

---

## 3. Rincian Endpoint & Contoh Penggunaan

### 3.1. Get Profile (Ambil Data Profil)

Mengambil data profil pengguna yang sedang login berdasarkan klaim token JWT.

- **HTTP Method:** `GET`
- **URL Path:** `/api/v1/profile`
- **Headers:**
  - `Authorization: Bearer <JWT_ACCESS_TOKEN>`

#### Parameter Request
*Tidak memerlukan query parameter maupun request body.*

#### Response Sukses (`200 OK`)
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "name": "Fajri Ramadhan",
    "email": "fajri@example.com",
    "phone": "+6281234567890",
    "photo": "/uploads/profile/550e8400-e29b-41d4-a716-446655440000.jpg",
    "role": "user",
    "is_active": true,
    "is_2fa_enabled": false,
    "notification_secret": "whsec_550e8400e29b41d4a716446655440000",
    "created_at": "2026-10-16T10:00:00Z",
    "updated_at": "2026-10-16T10:00:00Z"
  },
  "error": null,
  "meta": null
}
```

#### Contoh Request cURL
```bash
curl -X GET "http://localhost:8080/api/v1/profile" \
  -H "Authorization: Bearer <JWT_ACCESS_TOKEN>"
```

---

### 3.2. Update Profile (Perbarui Data Profil)

Memperbarui informasi profil pengguna. Semua field bersifat opsional; hanya field yang dikirimkan yang akan diperbarui. Jika `email` atau `phone` diubah, sistem akan memastikan nilai baru tersebut belum digunakan oleh pengguna lain.

- **HTTP Method:** `PUT`
- **URL Path:** `/api/v1/profile`
- **Headers:**
  - `Authorization: Bearer <JWT_ACCESS_TOKEN>`
  - `Content-Type: application/json`

#### Request Body
| Field | Tipe | Wajib | Validasi / Aturan | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| `name` | string | Tidak | `min=2, max=100` | Nama lengkap baru |
| `email` | string | Tidak | `email`, Unik | Alamat email baru (dicek keunikannya) |
| `phone` | string | Tidak | `min=8, max=20`, Unik | Nomor telepon baru (dicek keunikannya) |
| `is_2fa_enabled` | boolean (pointer) | Tidak | `true` atau `false` | Status aktivasi 2FA |

#### Contoh Request Body
```json
{
  "name": "Fajri Ramadhan Updated",
  "email": "fajri.new@example.com",
  "phone": "+6281298765432",
  "is_2fa_enabled": true
}
```

#### Response Sukses (`200 OK`)
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "name": "Fajri Ramadhan Updated",
    "email": "fajri.new@example.com",
    "phone": "+6281298765432",
    "photo": "http://localhost:8080/uploads/profile/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg",
    "role": "user",
    "is_active": true,
    "is_2fa_enabled": true,
    "created_at": "2026-10-16T10:00:00Z",
    "updated_at": "2026-10-16T10:45:00Z"
  },
  "error": null,
  "meta": null
}
```

#### Contoh Request cURL
```bash
curl -X PUT "http://localhost:8080/api/v1/profile" \
  -H "Authorization: Bearer <JWT_ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Fajri Ramadhan Updated",
    "phone": "+6281298765432"
  }'
```

---

### 3.3. Toggle 2FA (Aktifkan / Nonaktifkan 2FA)

Endpoint khusus untuk mengubah status Two-Factor Authentication via email tanpa perlu mengirimkan seluruh data profil.

- **HTTP Method:** `PATCH`
- **URL Path:** `/api/v1/profile/2fa`
- **Headers:**
  - `Authorization: Bearer <JWT_ACCESS_TOKEN>`
  - `Content-Type: application/json`

#### Request Body
| Field | Tipe | Wajib | Keterangan |
| :--- | :--- | :--- | :--- |
| `enable` | boolean | Ya | `true` untuk mengaktifkan, `false` untuk menonaktifkan |

#### Contoh Request Body
```json
{
  "enable": true
}
```

#### Response Sukses Saat Mengaktifkan (`200 OK`)
```json
{
  "success": true,
  "message": "Two-factor authentication enabled successfully",
  "data": {
    "id": 1,
    "name": "Fajri Ramadhan",
    "email": "fajri@example.com",
    "phone": "+6281234567890",
    "role": "user",
    "is_active": true,
    "is_2fa_enabled": true,
    "created_at": "2026-10-16T10:00:00Z",
    "updated_at": "2026-10-16T10:50:00Z"
  },
  "error": null,
  "meta": null
}
```

#### Response Sukses Saat Menonaktifkan (`200 OK`)
```json
{
  "success": true,
  "message": "Two-factor authentication disabled successfully",
  "data": {
    "id": 1,
    "name": "Fajri Ramadhan",
    "email": "fajri@example.com",
    "phone": "+6281234567890",
    "role": "user",
    "is_active": true,
    "is_2fa_enabled": false,
    "created_at": "2026-10-16T10:00:00Z",
    "updated_at": "2026-10-16T10:52:00Z"
  },
  "error": null,
  "meta": null
}
```

#### Contoh Request cURL
```bash
curl -X PATCH "http://localhost:8080/api/v1/profile/2fa" \
  -H "Authorization: Bearer <JWT_ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "enable": true
  }'
```

---

### 3.4. Update Profile Photo (Upload Foto Profil)

Mengunggah dan memperbarui foto profil / avatar pengguna. File gambar akan disimpan ke folder `./uploads/profile/` dengan nama berkas unik berbasis UUID.

- **HTTP Method:** `PATCH`
- **URL Path:** `/api/v1/profile/photo`
- **Headers:**
  - `Authorization: Bearer <JWT_ACCESS_TOKEN>`
  - `Content-Type: multipart/form-data`

#### Form Data Parameters
| Field | Tipe | Wajib | Keterangan |
| :--- | :--- | :--- | :--- |
| `photo` | file | Ya | File gambar avatar (format yang didukung: PNG, JPG, JPEG, WEBP; ukuran maks: 5MB) |

#### Response Sukses (`200 OK`)
```json
{
  "success": true,
  "message": "Profile photo updated successfully",
  "data": {
    "id": 1,
    "name": "Fajri Ramadhan",
    "email": "fajri@example.com",
    "phone": "+6281234567890",
    "photo": "http://localhost:8080/uploads/profile/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg",
    "role": "user",
    "is_active": true,
    "is_2fa_enabled": false,
    "created_at": "2026-10-16T10:00:00Z",
    "updated_at": "2026-10-16T17:15:00Z"
  },
  "error": null,
  "meta": null
}
```

#### Contoh Request cURL
```bash
curl -X PATCH "http://localhost:8080/api/v1/profile/photo" \
  -H "Authorization: Bearer <JWT_ACCESS_TOKEN>" \
  -F "photo=@/path/to/my_avatar.jpg"
```

---

### 3.5. Change Password (Ganti Password)

Mengganti password akun saat ini dengan memverifikasi password lama. Setelah password berhasil diperbarui di database, sistem mengirimkan email konfirmasi keamanan secara otomatis.

- **HTTP Method:** `PATCH`
- **URL Path:** `/api/v1/profile/change-password`
- **Headers:**
  - `Authorization: Bearer <JWT_ACCESS_TOKEN>`
  - `Content-Type: application/json`

#### Request Body
| Field | Tipe | Wajib | Validasi / Aturan | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| `old_password` | string | Ya | - | Password yang sedang aktif |
| `new_password` | string | Ya | `min=6` | Password baru |
| `confirm_password` | string | Ya | `eqfield=NewPassword` | Konfirmasi password baru (harus sama persis dengan `new_password`) |

#### Contoh Request Body
```json
{
  "old_password": "SecurePassword123",
  "new_password": "NewSecretPassword2026!",
  "confirm_password": "NewSecretPassword2026!"
}
```

#### Response Sukses (`200 OK`)
```json
{
  "success": true,
  "message": "Password changed successfully",
  "data": null,
  "error": null,
  "meta": null
}
```

#### Contoh Request cURL
```bash
curl -X PATCH "http://localhost:8080/api/v1/profile/change-password" \
  -H "Authorization: Bearer <JWT_ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "SecurePassword123",
    "new_password": "NewSecretPassword2026!",
    "confirm_password": "NewSecretPassword2026!"
  }'
```

---

### 3.6. Regenerate Notification Secret (Buat Ulang Secret Notifikasi)

Membuat ulang secret key khusus untuk endpoint webhook monitoring notifikasi Android (`whsec_...`). Gunakan endpoint ini jika secret lama Anda bocor atau ingin di-reset.

- **HTTP Method:** `POST` atau `PATCH`
- **URL Path:** `/api/v1/profile/notification-secret/regenerate`
- **Headers:**
  - `Authorization: Bearer <JWT_ACCESS_TOKEN>`

#### Parameter Request
*Tidak memerlukan query parameter maupun request body.*

#### Response Sukses (`200 OK`)
```json
{
  "success": true,
  "message": "Notification webhook secret regenerated successfully",
  "data": {
    "id": 1,
    "name": "Fajri Ramadhan",
    "email": "fajri@example.com",
    "phone": "+6281234567890",
    "photo": "http://localhost:8080/uploads/profile/550e8400-e29b-41d4-a716-446655440000.jpg",
    "role": "user",
    "is_active": true,
    "is_2fa_enabled": false,
    "notification_secret": "whsec_b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7",
    "created_at": "2026-10-16T10:00:00Z",
    "updated_at": "2026-10-16T22:05:00Z"
  },
  "error": null,
  "meta": null
}
```

#### Contoh Request cURL
```bash
curl -X POST "http://localhost:8080/api/v1/profile/notification-secret/regenerate" \
  -H "Authorization: Bearer <JWT_ACCESS_TOKEN>"
```

---

## 4. Penanganan Error & Kode Status HTTP

| Status Code | Kondisi Penyebab | Contoh Format Respon |
| :--- | :--- | :--- |
| **`401 Unauthorized`** | Header `Authorization` tidak disertakan, token kedaluwarsa, atau tanda tangan token tidak valid. | ```json { "success": false, "message": "Invalid or expired token", "data": null, "error": null, "meta": null } ``` |
| **`400 Bad Request`** | Password lama salah (`incorrect current password`), email/no. HP sudah terdaftar (`email already taken`), atau parameter tidak valid. | ```json { "success": false, "message": "incorrect current password", "data": null, "error": null, "meta": null } ``` |
| **`422 Unprocessable Entity`** | Gagal validasi DTO (misal: panjang nama tidak sesuai, format email salah, atau konfirmasi password tidak cocok). | ```json { "success": false, "message": "Validation failed", "data": null, "error": { "confirm_password": "confirm_password must be equal to new_password", "new_password": "new_password must be at least 6 characters" }, "meta": null } ``` |
| **`404 Not Found`** | ID Pengguna pada token tidak ditemukan di database (misal record telah dihapus secara soft-delete). | ```json { "success": false, "message": "user not found", "data": null, "error": null, "meta": null } ``` |

---

## 5. Contoh Integrasi Frontend (JavaScript / Axios)

Berikut adalah modul utilitas JavaScript untuk mempermudah integrasi fitur Profile di aplikasi frontend (React / Vue / Next.js / React Native):

```javascript
import axios from "axios";

const API_BASE_URL = "http://localhost:8080/api/v1";

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Interceptor menyematkan Bearer Token dari storage
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem("access_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 1. Mengambil data profil user saat ini
export async function getProfile() {
  try {
    const response = await apiClient.get("/profile");
    return response.data.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
}

// 2. Memperbarui informasi profil
export async function updateProfile({ name, email, phone }) {
  try {
    const response = await apiClient.put("/profile", { name, email, phone });
    return response.data.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
}

// 3. Mengaktifkan / Menonaktifkan 2FA
export async function setTwoFactorAuth(enable) {
  try {
    const response = await apiClient.patch("/profile/2fa", { enable });
    return response.data.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
}

// 4. Mengunggah foto profil (multipart/form-data)
export async function updateProfilePhoto(file) {
  const formData = new FormData();
  formData.append("photo", file);

  try {
    const response = await apiClient.patch("/profile/photo", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return response.data.data;
  } catch (error) {
    throw error.response?.data || error.message;
  }
}

// 5. Mengubah password akun
export async function changePassword(oldPassword, newPassword, confirmPassword) {
  try {
    const response = await apiClient.patch("/profile/change-password", {
      old_password: oldPassword,
      new_password: newPassword,
      confirm_password: confirmPassword,
    });
    return response.data.message;
  } catch (error) {
    throw error.response?.data || error.message;
  }
}
```
