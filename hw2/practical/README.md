<div dir="rtl" align="right">

<h1 align="right" dir="rtl">تمرین عملی — احراز هویت با پراکسی</h1>

<h2 align="right" dir="rtl">سؤال ۱ — اجرای سرویس فقط روی شبکهٔ داخلی</h2>

<p align="right" dir="rtl">سرویس را طوری بالا می‌آوریم که از اینترنت مستقیم به آن نرسد و فقط <code dir="ltr">Reverse Proxy</code>، که روی گیت‌وی است، بتواند با آن حرف بزند.</p>

<ul dir="rtl" align="right">
<li align="right" dir="rtl"><strong>آدرس اتصال (<code dir="ltr">Bind Address</code>):</strong> سرویس به‌جای گوش دادن روی همهٔ اینترفیس‌ها (<code dir="ltr">0.0.0.0</code>) فقط روی <code dir="ltr">IP</code> داخلی شبکه گوش می‌دهد، مثل <code dir="ltr">192.168.x.x</code> یا <code dir="ltr">10.x.x.x</code>. اگر پروکسی و سرویس روی یک ماشین باشند، می‌توان آن را روی <code dir="ltr">127.0.0.1</code> هم بست. در این تمرین پروکسی روی گیت‌وی است و سرویس روی ماشین دیگری در شبکهٔ داخلی، پس بستن سرویس فقط روی <code dir="ltr">127.0.0.1</code> درست نیست، چون پروکسی به آن نمی‌رسد.</li>
<li align="right" dir="rtl"><strong>فایروال سیستم‌عامل:</strong> با ابزاری مثل <code dir="ltr">ufw</code>، <code dir="ltr">iptables</code> یا <code dir="ltr">nftables</code> پورت سرویس (مثلاً <code dir="ltr">8081</code>) را طوری محدود می‌کنیم که فقط از رنج شبکهٔ داخلی یا از <code dir="ltr">IP</code> خود پروکسی قبول شود و بقیه <code dir="ltr">DROP</code> یا <code dir="ltr">REJECT</code> شوند.</li>
<li align="right" dir="rtl"><strong>جدا کردن شبکه:</strong> سرور را در یک زیرشبکهٔ خصوصی می‌گذاریم، به آن <code dir="ltr">IP</code> عمومی نمی‌دهیم و روی روتر مرزی هم <code dir="ltr">Port Forwarding</code> برای آن تنظیم نمی‌کنیم. فقط <code dir="ltr">Reverse Proxy</code> از بیرون دیده می‌شود و درخواست‌های پورت <code dir="ltr">80</code> را به پورت داخلی می‌فرستد.</li>
<li align="right" dir="rtl"><strong>شبکهٔ داکر:</strong> اگر سرویس داخل داکر باشد، پورت آن را با <code dir="ltr">-p</code> به بیرون <code dir="ltr">publish</code> نمی‌کنیم. پروکسی و سرویس را داخل یک شبکهٔ خصوصی داکر می‌گذاریم تا فقط پروکسی، با نام کانتینر، به سرویس برسد.</li>
</ul>

<h2 align="right" dir="rtl">سؤال ۲ — فلسفهٔ TLS Termination</h2>

<p align="right" dir="rtl">فلسفهٔ <code dir="ltr">TLS Termination</code> این است که رمزنگاری <code dir="ltr">HTTPS</code> همان اول، روی <code dir="ltr">Reverse Proxy</code>، تمام شود.</p>

<p align="right" dir="rtl">کاربر از اینترنت با <code dir="ltr">HTTPS</code> به پروکسی وصل می‌شود. پروکسی گواهی را دارد، ترافیک را باز می‌کند، درخواست را می‌بیند و مثلاً <code dir="ltr">JWT</code> را چک می‌کند. بعد همان درخواست را داخل شبکهٔ داخلی، که امن فرض می‌شود، به‌صورت <code dir="ltr">HTTP</code> معمولی به سرویس قدیمی می‌فرستد.</p>

<p align="right" dir="rtl">سرویس <code dir="ltr">Legacy</code> لازم نیست خودش <code dir="ltr">HTTPS</code> را پیاده کند یا گواهی داشته باشد؛ خیلی از سرویس‌های قدیمی اصلاً این قابلیت را ندارند و لازم نیست کدشان عوض شود. گواهی فقط یک جا، یعنی روی پروکسی، نصب و تمدید می‌شود و بار رمزنگاری هم از روی سرویس قدیمی برداشته می‌شود. پروکسی هم بعد از باز کردن <code dir="ltr">TLS</code> می‌تواند محتوا را بررسی کند، الگوی مخرب را رد کند و هدر احراز هویت را به درخواست اضافه کند.</p>

<h2 align="right" dir="rtl">سؤال ۳ — پروکسی احراز هویت JWT</h2>

<p align="right" dir="rtl">سرویس قدیمی همان ایمیج <code dir="ltr">ghcr.io/mccutchen/go-httpbin</code> است و داخل شبکهٔ داکر روی پورت <code dir="ltr">8081</code> گوش می‌دهد. پورت آن به ماشین میزبان منتشر نشده، پس از بیرون مستقیم به آن نمی‌رسد و فقط پروکسی، با نام <code dir="ltr">legacy</code>، به آن وصل می‌شود.</p>

<p align="right" dir="rtl">پروکسی با زبان <code dir="ltr">Go</code> نوشته شده و تنها سرویسی است که از بیرون دیده می‌شود. روی میزبان پورت <code dir="ltr">80</code> را می‌گیرد و درخواست را، بعد از بررسی توکن، به <code dir="ltr">http://legacy:8081</code> می‌فرستد. مسیر <code dir="ltr">POST /login</code> مال خود پروکسی است و به سرویس قدیمی نمی‌رود. با نام کاربری و رمز درست یک <code dir="ltr">JWT</code> با الگوریتم <code dir="ltr">HS256</code> برمی‌گرداند که یک ساعت اعتبار دارد.</p>

<p align="right" dir="rtl">بقیهٔ مسیرها بدون هدر <code dir="ltr">Authorization: Bearer</code> و توکن معتبر پاسخ <code dir="ltr">401</code> می‌گیرند و اصلاً به <code dir="ltr">httpbin</code> نمی‌رسند. توکن خراب، منقضی، یا امضاشده با کلید دیگر هم رد می‌شود. اگر توکن درست باشد، پروکسی نام کاربر را از claim به نام <code dir="ltr">sub</code> می‌خواند، هدر <code dir="ltr">X-Authenticated-User</code> را با همان نام پر می‌کند و درخواست را به سرویس داخلی می‌فرستد.</p>

<p align="right" dir="rtl">مقادیر پیش‌فرض فقط برای اجرای محلی تمرین هستند: نام کاربری <code dir="ltr">student</code>، رمز <code dir="ltr">student</code> و کلید <code dir="ltr">dev-only-change-this-secret</code>. برای اجرا از پوشهٔ <code dir="ltr">hw2/practical</code> این دستور را می‌زنیم. پروکسی روی پورت <code dir="ltr">80</code> بالا می‌آید. اگر آن پورت روی ماشین اشغال باشد، می‌شود پورت میزبان را عوض کرد، مثلاً <code dir="ltr">PROXY_PORT=9080 docker compose up --build</code>.</p>

</div>

```bash
docker compose up --build
```

<div dir="rtl" align="right">

<p align="right" dir="rtl">گرفتن توکن و امتحان کردن پروکسی:</p>

</div>

```bash
curl -s -X POST http://localhost/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"student","password":"student"}'

curl -i http://localhost/get

curl -i http://localhost/get \
  -H "Authorization: Bearer TOKEN"
```

<div dir="rtl" align="right">

<p align="right" dir="rtl">درخواست بدون توکن باید <code dir="ltr">401</code> بدهد. درخواست با توکن باید جواب <code dir="ltr">httpbin</code> را برگرداند. اتصال مستقیم به پورت <code dir="ltr">8081</code> روی خود ماشین نباید برقرار شود، چون آن پورت منتشر نشده است.</p>

</div>
