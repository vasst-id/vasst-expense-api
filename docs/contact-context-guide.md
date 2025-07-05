# 📋 **Contact Context Template Guide**

## **Overview**

The Contact Context Template is a streamlined JSON structure that maintains essential conversational context for each contact in the VASST Communication Agent. This template gets updated dynamically as messages are exchanged, with all information extracted directly from conversations and messages.

## **Template Location**
- **File**: `docs/contact-context-template.json`
- **Usage**: Initialize new contacts or update existing contact context

## **🎯 AI Filling Strategy**

**Core Principle**: All information must be extracted from actual conversations and messages. Never assume or fabricate data.

### **Information Sources Priority**:
1. **Direct statements** from the contact
2. **Contextual clues** from conversation flow
3. **Behavioral patterns** observed in messages
4. **Implicit information** from conversation context

---

## **Structure Breakdown**

### **🧑 Basic Contact Information**
```json
{
  "contact_id": "uuid-generated-by-system",
  "organization_id": "uuid-from-system",
  "name": "Budi Santoso",
  "phone_number": "+628123456789",
  "email": "budi.santoso@email.com",
  "address": "Jl. Sudirman No. 123, Jakarta Selatan",
  "salutation": "Pak"
}
```

**AI Filling Guidelines:**
- **name**: Extract from order form in the message
- **phone_number**: System-provided from messaging platform
- **email**: Only when explicitly mentioned in conversation
- **address**: Extract from  order form in the message
- **salutation**: Use appropriate Indonesian honorifics (Pak/Bu/Mas/Mbak, default using Kak)

**Examples from conversation:**
"Untuk pemesanan, boleh bantu diisi ya kak ☺️

Order form:
1. Nama: sonia
2. Alamat Pengiriman: Jln bambu betung 3, Duri Kosambi (no A6/6 pagar hitam pencet bell)
3. Nomor Telepon: 081231231231
4. Pesanan: 3 pack telur kampung omega, 1 pack telur negeri omega
5. Catatan Tambahan (jika ada):
6. Referral code:"

### **👤 Customer Information**
```json
{
  "customer_info": {
    "type": "lead",
    "language": "id",
    "tags": ["sensitif_harga"],
    "address": "Jl. Sudirman No. 123, Jakarta Barat",
    "area": "Green Lake City",
    "address_note": "Komplek Graha Indah, Blok C No. 15",
    "delivery_fee": 15000,
    "favorite_product": ["telur kampung omega pack", "telur negeri omega pack"],
    "orders_count": 3,
    "last_order_date": "2025-01-15T14:30:00Z",
    "last_order_amount": 76000,
    "last_order_product": ["2 telur kampung omega pack"],
    "last_order_status": "dikirim",
    "satisfaction_score": 8,
    "response_style": "friendly",
    "usual_day_order": "jumat"
  }
}
```

**AI Filling Guidelines:**

#### **type**: Customer classification
- **"lead"**: First-time inquirer, no purchase history
- **"new"**: Made 1-2 purchases, still getting familiar
- **"repeat"**: Regular customer with 3+ orders

**Examples:**
- Lead: "Halo, saya mau tanya tentang telur"
- New: "Saya kemarin sudah pesen, mau pesan lagi"
- Repeat: "Seperti biasa ya, pesanan yang sama"

#### **language**: Communication language
- **"id"**: Bahasa Indonesia (default)
- **"en"**: English (when customer uses English)

#### **tags**: Customer characteristics from conversation
**Examples:**
- "potensial_besar" - Shows strong buying intent
- "butuh_follow_up" - Asks to be contacted later
- "sensitif_harga" - Always asks about discounts
- "suka_promo" - Interested in promotions
- "customer_vip" - High-value customer
- "komplain_sering" - Has frequent complaints

#### **address/area**: Location information
**Extract from order form:**
"Untuk pemesanan, boleh bantu diisi ya kak ☺️

Order form:
1. Nama: sonia
2. Alamat Pengiriman: Jln bambu betung 3, Duri Kosambi (no A6/6 pagar hitam pencet bell)
3. Nomor Telepon: 081231231231
4. Pesanan: 3 pack telur kampung omega, 1 pack telur negeri omega
5. Catatan Tambahan (jika ada):
6. Referral code:"

#### **address_note**: Special delivery instructions
**Extract from order form (Catatan tambahan):**
"Untuk pemesanan, boleh bantu diisi ya kak ☺️

Order form:
1. Nama: sonia
2. Alamat Pengiriman: Jln bambu betung 3, Duri Kosambi (no A6/6 pagar hitam pencet bell)
3. Nomor Telepon: 081231231231
4. Pesanan: 3 pack telur kampung omega, 1 pack telur negeri omega
5. Catatan Tambahan (jika ada):
6. Referral code:"

#### **delivery_fee**: Shipping cost
**Extract from order confirmation:**

#### **favorite_product**: Preferred items
**Track from repeat orders:**
- "Pesan lagi telur kampungnya 2 pack ya"

#### **orders_count**: Number of completed orders
**Count from conversation history:**
- Increment when order is completed
- Reset if customer explicitly mentions it's their first order

#### **last_order_date/amount/product**: Recent order details
**Update with each new order:**
- Date: When order was placed
- Amount: Total value mentioned
- Product: Items ordered

#### **last_order_status**: Current order state
**Values:**
- "menunggu_konfirmasi" - Order placed, awaiting confirmation
- "menunggu_pembayaran" - awaiting payment
- "diproses" - Being prepared
- "dikirim" - In delivery
- "selesai" - Completed
- "dibatalkan" - Cancelled
- "dikembalikan" - Returned

#### **satisfaction_score**: Customer happiness (1-10)
**Gauge from conversation tone:**
- 1-3: "Kecewa banget", "Tidak sesuai ekspektasi"
- 4-6: "Ya lumayan", "Biasa aja"
- 7-8: "Bagus nih", "Suka banget"
- 9-10: "Luar biasa!", "Perfect banget!"

#### **response_style**: Communication preference
**Values:**
- "formal" - Uses formal language, "Bapak/Ibu"
- "santai" - Casual tone, uses "kamu/aku"
- "singkat" - Prefers brief responses
- "detail" - Likes comprehensive information

#### **usual_day_order**: Preferred ordering day
**Track patterns:**
- "Biasanya saya pesan hari Jumat"
- "Kalau weekend ada promo?"
- Notice ordering patterns

### **💬 Active Context**
```json
{
  "active_context": {
    "current_topic": "menanyakan harga telur negeri",
    "description": "Customer bertanya detail harga dan ingin tahu promo yang tersedia untuk pembelian serta ongkir"
  }
}
```

**AI Filling Guidelines:**

#### **current_topic**: What's being discussed now
**Examples:**
- "menanyakan harga produk"
- "komplain tentang pengiriman"
- "komplain tentang kualitas produk"
- "konfirmasi alamat pengiriman"
- "tanya status pesanan"

#### **description**: Detailed context
**Summarize the conversation flow:**
- "Customer baru pertama kali bertanya, tertarik dengan telur negeri"
- "Pelanggan lama, mau repeat order tapi alamat berubah"
- "Ada masalah dengan pesanan kemarin, minta solusi"

### **📊 Session Summary**
```json
{
  "session_summary": {
    "started_at": "2024-01-15T10:30:00Z",
    "messages_count": 8,
    "summary": "Customer baru menanyakan produk kopi, tertarik dengan telur negeri, minta info harga untuk pembelian 1 tray",
    "sentiment": "positif",
    "needs_human": false,
    "last_question": "Kalau beli 2 pack dapat diskon berapa persen?"
  }
}
```

**AI Filling Guidelines:**

#### **started_at**: Session start timestamp
- Set when first message in current session received

#### **messages_count**: Number of messages exchanged
- Count both customer and agent messages in current session

#### **summary**: Conversation overview
**Write in Bahasa Indonesia, include:**
- Customer type (new/existing)
- Main topic discussed
- Key decisions or outcomes
- Next steps if any

**Examples:**
- "Pelanggan lama order ulang menu favorit, alamat sama, minta dikirim besok pagi"
- "Customer baru tanya-tanya produk, belum decide, mau dipikir dulu"
- "Komplain pesanan terlambat, sudah dikasih solusi, customer puas"

#### **sentiment**: Overall conversation mood
**Values based on language tone:**
- **"positif"**: "Terima kasih", "Bagus banget", "Suka"
- **"negatif"**: "Kecewa", "Tidak puas", "Bermasalah"
- **"netral"**: Normal inquiry, no strong emotions

#### **needs_human**: Requires human intervention
**Set to true when:**
- Complex complaints
- Special requests beyond AI capability
- Customer specifically asks for human agent
- Escalation required

**Examples:**
- "Saya mau bicara sama manajer / admin manusia"
- "Ini masalah serius, tolong hubungi atasan"
- Technical issues beyond AI knowledge

#### **last_question**: Customer's most recent question
**Extract exact question or paraphrase:**
- "Kalau beli 2kg dapat diskon berapa persen?"
- "Kapan pesanan saya dikirim?"
- "Ada jual produk lain tidak?"

### **🧠 Memory System**
```json
{
  "memory": {
    "important_facts": [
      "Untuk mpasi, harus pilih telur yang segar",
      "Lebih suka dikirim pagi hari sebelum jam 10",
      "Pelanggan hasil referral customer Ayu"
    ],
    "previous_issues": [
      "Ada 3 telur busuk di pemesanan tanggal 23 Juni 2025",
      "Pengiriman salah alamat untuk pemesanan tanggal 14 Maret 2025"
    ]
  }
}
```

**AI Filling Guidelines:**

#### **important_facts**: Key information to remember
**Include:**
- Allergies or dietary restrictions
- Delivery preferences
- Business context (reseller, personal use)
- Special requirements
- Personal preferences

**Examples from conversation:**
- "Tolong kirim pagi ya, siang saya tidak ada di rumah"
- "Telur nya dipakai untuk mpasi"

#### **previous_issues**: Past problems
**Track for better service:**
- Delivery problems
- Product complaints
- Service issues
- Quality concerns

**Format**: "Brief description + date/period"
**Examples:**
- "Telur kemarin bau busuk, ada 3 butir"
- "Ada telur yang retak"
- "Pesanan belum sampai, sepertinya salah alamat"

### **⚙️ System Fields**
```json
{
  "system_fields": {
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:45:00Z",
    "last_message_at": "2024-01-15T11:45:00Z",
    "version": "1.0",
    "context_health": "active",
    "token_count": 100
  }
}
```

**AI Filling Guidelines:**
- **created_at**: Set once when contact is first created
- **updated_at**: Update every time context is modified
- **last_message_at**: Timestamp of most recent message
- **version**: Schema version (system managed)
- **context_health**: 
  - "active" - Recent activity (< 7 days)
  - "stale" - No activity 7-30 days
  - "archived" - No activity > 30 days
- **token_count**: Calculate total token

---

## **🤖 AI Context Filling Best Practices**

### **1. Information Extraction Rules**

#### **Direct Extraction** (Highest Priority)
```
Customer: "Nama saya Budi, alamat di Jl. Sudirman 123"
→ name: "Budi", address: "Jl. Sudirman 123"
```

#### **Contextual Inference** (Medium Priority)
```
Customer: "Saya dari Jakarta Selatan"
→ area: "Jakarta Selatan"
```

#### **Behavioral Pattern** (Lower Priority)
```
Customer orders every Friday for 3 weeks
→ usual_day_order: "jumat"
```

### **2. Language and Tone Analysis**

#### **Formal Language Detection**
```
"Selamat pagi, Bu. Saya ingin menanyakan..."
→ response_style: "formal"
```

#### **Casual Language Detection**
```
"Halo, mau tanya dong..."
→ response_style: "santai"
```

#### **Sentiment Analysis Examples**
```
"Wah enak banget telurnya!" → sentiment: "positif"
"Kok lama banget sih kirimnya?" → sentiment: "negatif"
"Oke, terima kasih infonya" → sentiment: "netral"
```

### **3. Progressive Context Building**

#### **First Interaction**
```json
{
  "name": "Budi",
  "customer_info": {
    "type": "lead",
    "language": "id"
  },
  "session_summary": {
    "summary": "Customer baru menanyakan price list"
  }
}
```

#### **After Order Placement**
```json
{
  "customer_info": {
    "type": "new",
    "orders_count": 1,
    "last_order_date": "2024-01-15T14:30:00Z"
  }
}
```

#### **After Multiple Orders**
```json
{
  "customer_info": {
    "type": "repeat",
    "orders_count": 5,
    "favorite_product": ["telur negeri omega pack", "telur kampung 1 tray"]
  }
}
```

### **4. Error Prevention Guidelines**

#### **Never Assume Information**
❌ Wrong: Guessing customer location
✅ Correct: Only fill when explicitly mentioned

#### **Don't Fabricate Details**
❌ Wrong: Making up satisfaction scores
✅ Correct: Only rate based on actual feedback

#### **Validate Consistency**
❌ Wrong: Conflicting information across fields
✅ Correct: Cross-check all related fields

### **5. Update Triggers**

#### **Every Message Should Update:**
- `messages_count`
- `last_message_at` 
- `updated_at`
- `current_topic` (if topic changes)
- `sentiment` (if sentiment shifts)

#### **Specific Events Should Update:**
- **New order**: `orders_count`, `last_order_*` fields
- **Address mentioned**: `address`, `area`, `address_note`
- **Preference stated**: `favorite_product`, `response_style`
- **Problem reported**: Add to `previous_issues`
- **Important info**: Add to `important_facts`

---

## **🔧 Implementation Examples**

### **Example 1: New Customer Inquiry**

**Conversation:**
```
Customer: "Halo, mau tanya telur negeri omega nya harganya berapa ya?"
Agent: "Halo Kak! Telur negeri omega per pack harganya 30.000"
Customer: "Ada promo gak?"
```

**Context Update:**
```json
{
  "name": "-",
  "salutation": "Kak",
  "customer_info": {
    "type": "lead",
    "language": "id",
    "tags": ["suka_promo"]
  },
  "active_context": {
    "current_topic": "menanyakan harga telur negeri omega",
    "description": "Customer baru tanya harga, menanyakan promo di awal"
  },
  "session_summary": {
    "messages_count": 3,
    "summary": "Customer baru bernama menanyakan harga telur negeri omega, tanya promo di awal",
    "sentiment": "netral",
    "last_question": "Ada promo gak?"
  }
}
```

### **Example 2: Regular Customer Reorder**

**Conversation:**
```
Customer: "Halo kak, pesan telur nya lagi yah, 2 pack telur kampung omega"
Agent: "Siap Kak! kami proses segera"
Customer: "Ok. Tolong kirim besok pagi ya"
```

**Context Update:**
```json
{
  "customer_info": {
    "type": "repeat",
    "orders_count": 8,
    "last_order_date": "2024-01-15T14:30:00Z",
    "last_order_amount": 420000,
    "last_order_product": ["2 x telur kampung omega pack"],
    "last_order_status": "dikirim",
    "favorite_product": ["telur kampung omega pack"]
  },
  "active_context": {
    "current_topic": "repeat order produk favorit",
    "description": "Pelanggan tetap order ulang menu biasa, minta dikirim besok pagi"
  },
  "session_summary": {
    "summary": "Pak Budi repeat order menu favorit kopi 2 x telur kampung omega pack, minta kirim besok pagi",
    "sentiment": "positif"
  },
  "memory": {
    "important_facts": ["Lebih suka dikirim pagi hari"]
  }
}
```

### **Example 3: Customer Complaint**

**Conversation:**
```
Customer: "Kak, pesanan nya kok belum nyampe ya?"
Agent: "Halo kak, maaf kak, ada kendala di kurir kami. Kami cek dulu ya"
Customer: "Waduh, soalnya saya mau pake buat masak"
```

**Context Update:**
```json
{
  "active_context": {
    "current_topic": "komplain keterlambatan pengiriman",
    "description": "Customer komplain pesanan belum sampai, butuh untuk masak, perlu solusi cepat"
  },
  "session_summary": {
    "sentiment": "negatif",
    "needs_human": true,
    "summary": "Customer komplain pesanan terlambat, dibutuhkan untuk masak, perlu penanganan khusus"
  },
  "memory": {
    "previous_issues": ["Pesanan terlambat untuk order tanggal 24 Juni 2025"]
  }
}
```