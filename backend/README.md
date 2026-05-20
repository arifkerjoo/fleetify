
=======================================  MVP ENDPOINTS  ==============================================

[MASTER DATA]
GET /api/v1/vehicles       : Ambil daftar kendaraan untuk dipilih oleh SA saat membuat laporan.
GET /api/v1/items             : Ambil master part/jasa beserta harga terbaru.

[REPORT - CORE]
POST /api/v1/reports          : Buat laporan pemeliharaan baru (SA only. STATUS : PENDING_APPROVAL).  
    FORM : {
        kendaraan,
        odometer,
        keluhan,
        foto simulasi,
        list item (part/jasa + qty)
    }

GET /api/v1/reports            : Ambil semua laporan (SA & Approval). Filter opsional: status, tanggal
GET /api/v1/reports/:id        : Detail lengkap satu laporan (header + items + status history).

[FLOW - CORE]
PATCH /api/reports/:id/approve  : Approval menyetujui laporan. (STATUS : APPROVED)
PATCH /api/reports/:id/complete : SA menyelesaikan pekerjaan dengan upload foto bukti. (STATUS : COMPLETED)

========================================================================

OPSIONAL : 
GET /api/reports/export/csv     : Download semua laporan dalam format CSV (native JS frontend consumer).
    Isi {
        SA name
        nomor polisi
        status
        tanggal
        total estimasi
    }
