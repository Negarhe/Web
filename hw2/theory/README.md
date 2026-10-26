<div dir="rtl" align="right">

<h1 align="right" dir="rtl">تمرین HTTP</h1>

<h2 align="right" dir="rtl">سؤال ۱ — تفاوت نسخه‌های HTTP</h2>

<h3 align="right" dir="rtl">تفاوت‌های کلیدی</h3>

<ul dir="rtl" align="right">
<li align="right" dir="rtl"><strong><code dir="ltr">HTTP/0.9</code>:</strong> فقط متد <code dir="ltr">GET</code> • بدون هدر و وضعیت (<code dir="ltr">Status Code</code>) • فقط ارسال فایل‌های <code dir="ltr">HTML</code> • قطع اتصال بلافاصله پس از دریافت پاسخ</li>
<li align="right" dir="rtl"><strong><code dir="ltr">HTTP/1.0</code>:</strong> اضافه شدن متدهای <code dir="ltr">POST</code> و <code dir="ltr">HEAD</code> • اضافه شدن هدرها و <code dir="ltr">Status Code</code>ها • پشتیبانی از انواع فایل‌ها (<code dir="ltr">Content-Type</code>) • نیاز به ایجاد کانکشن جدید <code dir="ltr">TCP</code> برای هر درخواست (<code dir="ltr">Non-persistent</code>)</li>
<li align="right" dir="rtl"><strong><code dir="ltr">HTTP/1.1</code>:</strong> اتصال پایدار (<code dir="ltr">Persistent Connection / Keep-Alive</code>) • امکان ارسال متوالی درخواست‌ها (<code dir="ltr">Pipelining</code>) • اجباری شدن هدر <code dir="ltr">Host</code> • پشتیبانی از <code dir="ltr">Chunked Transfer</code> و <code dir="ltr">Range Requests</code> • مشکل مسدود شدن صف در لایه کاربرد (<code dir="ltr">Head-of-Line Blocking</code>)</li>
<li align="right" dir="rtl"><strong><code dir="ltr">HTTP/2</code>:</strong> قالب باینری (<code dir="ltr">Binary Framing</code>) به جای متنی • چندارسال‌گری یا <code dir="ltr">Multiplexing</code> روی یک کانکشن <code dir="ltr">TCP</code> • فشرده‌سازی هدر با الگوریتم <code dir="ltr">HPACK</code> • قابلیت <code dir="ltr">Server Push</code> و اولویت‌بندی استریم‌ها</li>
<li align="right" dir="rtl"><strong><code dir="ltr">HTTP/3</code>:</strong> استفاده از پروتکل <code dir="ltr">QUIC</code> مبتنی بر <code dir="ltr">UDP</code> به جای <code dir="ltr">TCP</code> • فشرده‌سازی هدر با <code dir="ltr">QPACK</code> • ادغام مراحل ایجاد اتصال و رمزنگاری با <code dir="ltr">TLS 1.3</code> • استقلال کامل استریم‌ها</li>
</ul>

<h3 align="right" dir="rtl">چرا نسخه ۳ از UDP استفاده می‌کند؟</h3>

<p align="right" dir="rtl">در <code dir="ltr">HTTP/2</code> اگرچه چندین استریم روی یک کانکشن <code dir="ltr">TCP</code> مالتی‌پلکس می‌شوند، اما چون <code dir="ltr">TCP</code> کل ترافیک را به عنوان یک جریان پیوسته بایت در نظر می‌گیرد، اگر حتی یک بسته گم شود، تحویل تمام داده‌های بعدی (حتی استریم‌های نامرتبط دیگر) تا زمان ارسال مجدد آن پکت متوقف می‌شود (<code dir="ltr">Transport-level Head-of-Line Blocking</code>).</p>

<p align="right" dir="rtl">هدف از مهاجرت به <code dir="ltr">UDP</code> در <code dir="ltr">HTTP/3</code> عدم اطمینان در انتقال نیست؛ بلکه هدف این است که از بستر سبک و بدون قفل <code dir="ltr">UDP</code> استفاده شود تا پروتکل <code dir="ltr">QUIC</code> بتواند در فضای کاربر (<code dir="ltr">User-space</code>) مدیریت بسته‌ها، رمزنگاری یکپارچه (ادغام لایه انتقال با <code dir="ltr">TLS 1.3</code> برای برقراری سریع‌تر اتصال / کاهش <code dir="ltr">RTT</code>) و مدیریت مستقل استریم‌ها را انجام دهد. همچنین پیاده‌سازی روی <code dir="ltr">UDP</code> از مشکل دستکاری یا مسدودسازی آپدیت‌های جدید <code dir="ltr">TCP</code> توسط تجهیزات میانی شبکه (<code dir="ltr">Middleboxes / Protocol Ossification</code>) جلوگیری کرده و قابلیت مهاجرت اتصال (<code dir="ltr">Connection Migration</code>) را در زمان تغییر <code dir="ltr">IP</code> کلاینت فراهم می‌کند.</p>

<h2 align="right" dir="rtl">سؤال ۲ — ارتباط Stateful و Stateless</h2>

<h3 align="right" dir="rtl">آ)</h3>

<p align="right" dir="rtl">در ارتباط <code dir="ltr">Stateless</code> (بی‌حالت)، هیچ‌کدام از طرفین اطلاعاتی از تعاملات قبلی ذخیره نمی‌کنند و هر درخواست به‌صورت کاملاً مستقل و شامل تمام داده‌های لازم پردازش می‌شود. این مدل مقیاس‌پذیری <code dir="ltr">(Scalability)</code> بالایی دارد و مدیریت سرور را ساده می‌کند. در مقابل، در ارتباط <code dir="ltr">Stateful</code> (حالت‌دار)، وضعیت نشست یا همان <code dir="ltr">Context</code> گفت‌وگو بین درخواست‌ها حفظ می‌شود (مانند به‌خاطر سپردن وضعیت لاگین یا محتویات سبد خرید)؛ با اینکه این روش تجربه کاربری پیوسته‌ای ایجاد می‌کند، اما توزیع بار <code dir="ltr">(Load Balancing)</code> روی چند سرور و مدیریت خطای سرور را پیچیده می‌کند.</p>

<h3 align="right" dir="rtl">ب)</h3>

<p align="right" dir="rtl">پروتکل <code dir="ltr">HTTP</code> ذاتا یک پروتکل <code dir="ltr">Stateless</code> است. سرور پس از ارسال پاسخ، هیچ وضعیتی از کلاینت یا درخواست قبلی را در حافظه نگه نمی‌دارد. مشکل اصلی این ویژگی، ناتوانی در حفظ نشست‌های کاربری (مانند لاگین ماندن یا ثبت سبد خرید) در درخواست‌های متوالی است. برای حل این مشکل بدون تغییر ساختار پروتکل، از کوکی‌ها <code dir="ltr">(Cookies)</code> و نشست‌ها <code dir="ltr">(Sessions)</code> استفاده می‌شود:</p>

<p align="right" dir="rtl">سرور هنگام احراز هویت اولیه، یک شناسه یکتا <code dir="ltr">(Session ID)</code> ایجاد کرده و از طریق هدر <code dir="ltr">Set-Cookie</code> به کلاینت می‌فرستد. مرورگر این شناسه را ذخیره کرده و در تمام درخواست‌های بعدی، آن را داخل هدر <code dir="ltr">Cookie</code> به سرور بازمی‌گرداند. سرور با خواندن این شناسه، وضعیت مربوط به آن کاربر را از دیتابیس یا کش خود بازیابی می‌کند؛ بنابراین پروتکل <code dir="ltr">HTTP</code> همچنان <code dir="ltr">Stateless</code> باقی می‌ماند و صرفاً نقش حامل شناسه را بازی می‌کند.</p>

<h2 align="right" dir="rtl">سؤال ۳ — محدودیت شروع ارتباط در HTTP</h2>

<h3 align="right" dir="rtl">آ)</h3>

<p align="right" dir="rtl">در <code dir="ltr">HTTP</code> ارتباط همیشه با درخواست کلاینت شروع می‌شود و سرور فقط به همان درخواست جواب می‌دهد. سرور خودش نمی‌تواند اتصال را باز کند و داده را هل بدهد. در چت، اعلان یا <code dir="ltr">webhook</code> خبر سمت سرور اتفاق می‌افتد و کلاینت از قبل نمی‌داند کی باید درخواست بفرستد. اگر کلاینت مدام درخواست بفرستد، درخواست خالی زیاد می‌شود و هم پهنای باند هدر می‌رود هم سرور درگیر می‌شود. اگر هم دیر درخواست بفرستد، داده با تأخیر به دستش می‌رسد.</p>

<h3 align="right" dir="rtl">ب)</h3>

<p align="right" dir="rtl"><strong><span dir="ltr">Short Polling.</span></strong> کلاینت هر چند ثانیه یک‌بار یک درخواست <code dir="ltr">GET</code> می‌فرستد و می‌پرسد دادهٔ جدیدی هست یا نه. پیاده‌سازی‌اش ساده است، ولی اگر خبری نباشد کلی درخواست و جواب خالی رد و بدل می‌شود و بار سرور و ترافیک شبکه بالا می‌رود.</p>

<p align="right" dir="rtl"><strong><span dir="ltr">Long Polling.</span></strong> کلاینت باز هم درخواست می‌فرستد، ولی سرور تا وقتی دادهٔ جدید آماده شود یا زمان درخواست تمام شود جواب نمی‌دهد و اتصال را باز نگه می‌دارد. به محض رسیدن داده یا تمام شدن <code dir="ltr">timeout</code>، جواب را می‌فرستد و اتصال بسته می‌شود. کلاینت بلافاصله یک درخواست جدید می‌فرستد و دوباره منتظر می‌ماند. تأخیر این روش از <code dir="ltr">short polling</code> خیلی کمتر است، چون کلاینت لازم نیست هی در فاصله‌های کوتاه درخواست تکراری بزند.</p>

<p align="right" dir="rtl"><strong><span dir="ltr">Server-Sent Events (<code dir="ltr">SSE</code>).</span></strong> این روش یک اتصال <code dir="ltr">HTTP</code> باز از سرور به کلاینت نگه می‌دارد و داده را با نوع <code dir="ltr">text/event-stream</code> به‌صورت جریان متنی می‌فرستد. بعد از درخواست اول، سرور هر وقت اتفاقی افتاد همان اتصال را استفاده می‌کند و لازم نیست کلاینت دوباره درخواست بدهد. برای اعلان و مانیتورینگ مناسب است. مسیر فقط از سرور به کلاینت است و کلاینت روی همین کانال نمی‌تواند داده بفرستد.</p>

<p align="right" dir="rtl"><strong><span dir="ltr">WebSocket.</span></strong> اول یک درخواست عادی <code dir="ltr">HTTP</code> فرستاده می‌شود و با هدر <code dir="ltr">Upgrade</code> اتصال به <code dir="ltr">WebSocket</code> تبدیل می‌شود. بعد از این دست‌دادن، دیگر درخواست و جواب جدا نداریم و یک اتصال دوطرفه روی <code dir="ltr">TCP</code> باز می‌ماند. هم کلاینت هم سرور هر لحظه می‌توانند پیام بفرستند و هدر هر پیام هم کوچک است. برای چت و کارهایی که دو طرف باید سریع با هم حرف بزنند معمولاً از همین استفاده می‌شود.</p>

<p align="right" dir="rtl"><strong><span dir="ltr">Webhook.</span></strong> این راهکار برای ارتباط بین دو سرور است. سامانهٔ گیرنده یک <code dir="ltr">URL</code> ثبت می‌کند و سامانهٔ فرستنده، به محض رخ دادن رویداد، یک <code dir="ltr">POST</code> به همان آدرس می‌زند. یعنی کسی که باید داده را بگیرد خودش مثل یک سرور عمل می‌کند و لازم نیست مدام از طرف مقابل بپرسد خبری شده یا نه.</p>

<h2 align="right" dir="rtl">سؤال ۴ — اجزای URL</h2>

<p align="right" dir="rtl">آدرس:</p>

</div>

```
ss://asghar:1234!!@ss.myproxy.com:1234#shadowSocks1
```

<div dir="rtl" align="right">

<p align="right" dir="rtl">طرح یا همان <code dir="ltr">scheme</code> برابر <code dir="ltr">ss</code> است. این قسمت قبل از <code dir="ltr">://</code> آمده و نوع آدرس را مشخص می‌کند.</p>

<p align="right" dir="rtl">اطلاعات کاربری یا <code dir="ltr">userinfo</code> برابر <code dir="ltr">asghar:1234!!</code> است و بین <code dir="ltr">//</code> و <code dir="ltr">@</code> قرار دارد. نام کاربری <code dir="ltr">asghar</code> است و گذرواژه <code dir="ltr">1234!!</code>. هر دو علامت تعجب هم جزو رمز هستند.</p>

<p align="right" dir="rtl">میزبان <code dir="ltr">ss.myproxy.com</code> است و پورت <code dir="ltr">1234</code>. رمز همان <code dir="ltr">1234!!</code> قبل از <code dir="ltr">@</code> است و پورت عدد <code dir="ltr">1234</code> بعد از دو نقطه است.</p>

<p align="right" dir="rtl">مسیر (<code dir="ltr">path</code>) و کوئری (<code dir="ltr">query</code>) در این آدرس نیامده‌اند.</p>

<p align="right" dir="rtl">بعد از <code dir="ltr">#</code> مقدار <code dir="ltr">shadowSocks1</code> آمده که <code dir="ltr">fragment</code> است. این بخش با درخواست به سرور فرستاده نمی‌شود و سمت کلاینت می‌ماند. در کلاینت‌های پروکسی معمولاً همین را به‌عنوان اسم سرور نشان می‌دهند.</p>

<p align="right" dir="rtl">اگر <code dir="ltr">authority</code> را هم جدا کنیم، کل بخش <code dir="ltr">asghar:1234!!@ss.myproxy.com:1234</code> است؛ یعنی از بعد <code dir="ltr">//</code> تا قبل از <code dir="ltr">#</code>، که از اطلاعات کاربری، میزبان و پورت ساخته شده.</p>

<h2 align="right" dir="rtl">سؤال ۵ — کد وضعیت HTTP</h2>

<p align="right" dir="rtl"><strong>۱.</strong> کد <code dir="ltr">503 Service Unavailable</code> انتخاب بهتری است، چون دیتابیس خراب شده و سرور برای مدتی نمی‌تواند لاگین را انجام دهد. اگر این خطا داخل برنامه گرفته نشده باشد و سرور با یک خطای پیش‌بینی‌نشده مواجه شود، <code dir="ltr">500 Internal Server Error</code> هم پاسخ درستی است.</p>

<p align="right" dir="rtl"><strong>۲.</strong> کد <code dir="ltr">401 Unauthorized</code> مناسب است، چون هویت کاربر تأیید نشده. برای نام کاربری ناموجود هم همین کد را برمی‌گردانیم تا از روی وضعیت پاسخ معلوم نشود آن نام در سایت وجود دارد یا نه.</p>

<p align="right" dir="rtl"><strong>۳.</strong> کد <code dir="ltr">301 Moved Permanently</code> یعنی سایت برای همیشه به دامنهٔ جدید رفته است. در <code dir="ltr">REST API</code> کد <code dir="ltr">308 Permanent Redirect</code> انتخاب دقیق‌تری است، چون انتقال همیشگی است و متد درخواست، مثل <code dir="ltr">POST</code>، در آدرس جدید هم همان می‌ماند.</p>

<p align="right" dir="rtl"><strong>۴.</strong> کد <code dir="ltr">429 Too Many Requests</code> مخصوص وقتی است که کاربر از سقف مجاز درخواست رد شده باشد. اغلب کنارش هدر <code dir="ltr">Retry-After</code> هم می‌آید تا کلاینت بداند کی دوباره تلاش کند.</p>

<p align="right" dir="rtl"><strong>۵.</strong> چون تمدید موفق بوده و توکن جدید داخل هدر برگشته، کد <code dir="ltr">204 No Content</code> انتخاب دقیق‌تری است؛ یعنی کار انجام شده و بدنهٔ پاسخ خالی است. اگر علاوه بر هدر یک بدنهٔ <code dir="ltr">JSON</code> هم برگردانده شود، کد <code dir="ltr">200 OK</code> مناسب است.</p>

<p align="right" dir="rtl"><strong>۶.</strong> کد <code dir="ltr">403 Forbidden</code> پاسخ معمول این حالت است: سرور درخواست را فهمیده، ولی به این کلاینت سرویس نمی‌دهد. چون دلیل ممنوعیت تحریم است، کد <code dir="ltr">451 Unavailable For Legal Reasons</code> انتخاب دقیق‌تری است و مخصوص محدودیتی است که دلیل قانونی دارد.</p>

<h2 align="right" dir="rtl">سؤال ۶ — نوع پروکسی</h2>

<p align="right" dir="rtl"><strong>۱.</strong> نوعش <code dir="ltr">Forward Proxy</code> کش‌کننده است. سرور داخل دانشگاه، سمت کلاینت‌ها قرار دارد. سیستم‌ها برای گرفتن پکیج، به جای مخزن اصلی، درخواست را به این سرور می‌دهند. پروکسی بسته را از اینترنت می‌گیرد و ذخیره می‌کند و درخواست بعدی را از همان حافظه محلی جواب می‌دهد، پس اینترنت بین‌الملل کمتر مصرف می‌شود. برای <code dir="ltr">apt</code> ابزاری مثل <code dir="ltr">apt-cacher-ng</code> همین کار را می‌کند. اگر دانشگاه یک کپی محلی از مخزن داشته باشد، همان نقش را آینه (<code dir="ltr">Mirror</code>) هم بازی می‌کند.</p>

<p align="right" dir="rtl"><strong>۲.</strong> نوعش <code dir="ltr">Reverse Proxy</code> است که درخواست‌ها را بر اساس جغرافیا تقسیم می‌کند. لبه <code dir="ltr">CDN</code> هم همین نقش را دارد. پروکسی جلوی سرورهای اصلی است. درخواست از اینترنت می‌آید و پروکسی با توجه به محل کاربر آن را به نزدیک‌ترین مرکز داده می‌فرستد تا تاخیر کمتر شود و بار بین سرورها پخش شود.</p>

<p align="right" dir="rtl"><strong>۳.</strong> نوعش <code dir="ltr">Reverse Proxy</code> احراز هویت است. به این الگو <code dir="ltr">Ambassador</code> هم می‌گویند. پروکسی جلوی سرویس قدیمی قرار می‌گیرد و کلاینت با خود پروکسی ارتباط می‌گیرد. پروکسی هویت را بررسی می‌کند، مثلا با <code dir="ltr">JWT</code> یا نام کاربری و رمز. اگر کاربر مجاز باشد، درخواست به سرویس قدیمی فرستاده می‌شود. خود آن سرویس امکان لاگین ندارد.</p>

</div>
