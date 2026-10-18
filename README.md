<div dir="rtl" align="right">

<h1 align="right" dir="rtl">تمرین Git</h1>

<h2 align="right" dir="rtl">سؤال ۱ — <code dir="ltr">git merge</code> و <code dir="ltr">git rebase</code> (۳ نمره)</h2>

<p align="right" dir="rtl">هر دو تغییرات یک شاخه را وارد شاخۀ دیگر می‌کنند، ولی شکل تاریخچه فرق می‌کند.</p>

<p align="right" dir="rtl"><strong><code dir="ltr">git merge</code>:</strong> یک ادغام سه‌طرفه بین جد مشترک و نوک دو شاخه انجام می‌دهد و یک کامیت ادغام با دو والد می‌سازد (اگر <code dir="ltr">main</code> جلو نرفته باشد، فقط <span dir="ltr">fast-forward</span> می‌شود). کامیت‌های قبلی و هششان دست نمی‌خورند و تاریخچه غیرخطی می‌ماند.</p>

</div>

```text
A---B---C---F---M   (main)
     \         /
      D-------E     (feature)
```

<div dir="rtl" align="right">

<p align="right" dir="rtl"><strong><code dir="ltr">git rebase</code>:</strong> کامیت‌های شاخه را برمی‌دارد و یکی‌یکی روی نوک <code dir="ltr">main</code> دوباره اعمال می‌کند. کامیت‌های جدید (<code dir="ltr">D'</code>، <code dir="ltr">E'</code>) هش تازه دارند، یعنی تاریخچه بازنویسی می‌شود و خطی است. تعارض ممکن است برای هر کامیت جدا حل شود.</p>

</div>

```text
A---B---C---F           (main)
             \
              D'---E'   (feature)
```

<div dir="rtl" align="right">

<p align="right" dir="rtl"><strong>کاربرد <code dir="ltr">merge</code>:</strong> ادغام یک قابلیت تمام‌شده در شاخۀ مشترک، مثلاً ادغام <code dir="ltr">feature-login</code> در <code dir="ltr">main</code> با <span dir="ltr">pull request</span>. چون شاخه عمومی است، نباید کامیت‌ها بازنویسی شوند و کامیت ادغام هم قابل <code dir="ltr">revert</code> است.</p>

<p align="right" dir="rtl"><strong>کاربرد <code dir="ltr">rebase</code>:</strong> به‌روز کردن شاخۀ خصوصی خودم با آخرین <code dir="ltr">main</code> پیش از <span dir="ltr">pull request</span> (<code dir="ltr">git fetch && git rebase origin/main</code>). شاخه فقط مال من است، پس بازنویسی آسیبی ندارد و تاریخچه تمیز و خطی می‌شود.</p>

<p align="right" dir="rtl"><strong>قاعده:</strong> کامیت‌هایی را که روی شاخۀ مشترک پوش شده‌اند <code dir="ltr">rebase</code> نکنید.</p>

<h2 align="right" dir="rtl">سؤال ۲ — انتخاب بخشی از تغییرات (<span dir="ltr">Hunk</span>) (۳ نمره)</h2>

<p align="right" dir="rtl"><span dir="ltr">hunk</span> یک تکۀ پیوسته از تغییرات در خروجی <code dir="ltr">git diff</code> است که با سرآیندی مثل <code dir="ltr">@@ -14,3 +14,4 @@</code> شروع می‌شود (محل و تعداد خط‌ها در نسخۀ قدیم و جدید).</p>

</div>

```bash
git add -p README.md
git diff --staged
git commit -m "Clarify the setup steps"
```

<div dir="rtl" align="right">

<p align="right" dir="rtl">گیت هر <span dir="ltr">hunk</span> را نشان می‌دهد و می‌پرسد: <code dir="ltr">y</code> یعنی <span dir="ltr">stage</span> کن، <code dir="ltr">n</code> یعنی نکن، <code dir="ltr">s</code> یعنی به تکه‌های کوچک‌تر بشکن، <code dir="ltr">e</code> یعنی دستی ویرایش کن، <code dir="ltr">q</code> یعنی خروج. با <code dir="ltr">git diff --staged</code> می‌بینم دقیقاً چه چیزی کامیت می‌شود.</p>

<p align="right" dir="rtl"><strong>چرا مفید است:</strong> اگر یک فایل چند تغییر بی‌ربط داشته باشد (رفع خطا، تغییر ظاهری، کد اشکال‌زدایی)، می‌توان هر کدام را در کامیت جدا گذاشت و کد موقت را بیرون نگه داشت. هر کامیت یک کار انجام می‌دهد و بازبینی، <code dir="ltr">revert</code> و <code dir="ltr">bisect</code> ساده‌تر می‌شود.</p>

<h2 align="right" dir="rtl">سؤال ۳ — پیدا کردن منشأ خطا با <code dir="ltr">git bisect</code> (۳ نمره)</h2>

<p align="right" dir="rtl"><code dir="ltr">git bisect</code> بین یک کامیت سالم و یک کامیت معیوب جست‌وجوی دودویی می‌کند. در هر گام کامیت وسط را <span dir="ltr">checkout</span> می‌کند، ما آن را <code dir="ltr">good</code> یا <code dir="ltr">bad</code> علامت می‌زنیم و بازه نصف می‌شود. برای <span dir="ltr">N</span> کامیت حدود <span dir="ltr">log₂(N)</span> آزمایش کافی است.</p>

</div>

```bash
git bisect start
git bisect bad            # نسخۀ فعلی خراب است
git bisect good v1.8      # این نسخه سالم بود
# بعد از آزمایش هر کامیت:
git bisect good           # یا bad یا skip
git bisect reset          # پایان و برگشت به شاخه
```

<div dir="rtl" align="right">

<p align="right" dir="rtl"><strong>مثال:</strong> تعداد کالاها در نسخۀ فعلی غلط است ولی در <code dir="ltr">v1.8</code> (حدود ۱۸۰ کامیت قبل) درست بود. بعد از حدود ۸ گام گیت اعلام می‌کند <code dir="ltr">c8d41ab is the first bad commit</code> و با <code dir="ltr">git show c8d41ab</code> خط‌های مقصر دیده می‌شوند. با <code dir="ltr">git bisect run ./test.sh</code> هم می‌توان کل کار را خودکار کرد.</p>

<p align="right" dir="rtl"><strong>چرا مفید است:</strong> به‌جای بررسی صدها کامیت، با چند آزمایش و بدون حدس، دقیقاً کامیت مقصر (با <span dir="ltr">diff</span>، نویسنده و پیامش) پیدا می‌شود.</p>

<h2 align="right" dir="rtl">سؤال ۴ — بررسی تغییرات پیش از کامیت (۳ نمره)</h2>

<p align="right" dir="rtl">فرض: بخشی از تغییرات <code dir="ltr">README.md</code> (الف) <span dir="ltr">stage</span> شده، بعد فایل دوباره ویرایش شده (ب)، و <code dir="ltr">tmp/</code> <span dir="ltr">untracked</span> است.</p>

<h3 align="right" dir="rtl">۱. بررسی وضعیت</h3>

<ul dir="rtl" align="right">
<li align="right" dir="rtl"><code dir="ltr">git status</code>: نمای کلی؛ <code dir="ltr">README.md</code> هم در <span dir="ltr">staged</span> و هم در <span dir="ltr">unstaged</span> دیده می‌شود و <code dir="ltr">tmp/</code> در <span dir="ltr">Untracked</span>.</li>
<li align="right" dir="rtl"><code dir="ltr">git diff</code>: تغییرات <span dir="ltr">unstaged</span> (فقط ب).</li>
<li align="right" dir="rtl"><code dir="ltr">git diff --staged</code>: تغییرات <span dir="ltr">staged</span> که وارد کامیت بعدی می‌شوند (فقط الف).</li>
</ul>

<p align="right" dir="rtl">تغییر <span dir="ltr">staged</span> در <span dir="ltr">index</span> ثبت شده و کامیت می‌شود؛ تغییر <span dir="ltr">unstaged</span> فقط روی دیسک است تا وقتی <code dir="ltr">git add</code> شود.</p>

<h3 align="right" dir="rtl">۲. کامیت بدون <code dir="ltr">-a</code></h3>

<p align="right" dir="rtl">با <code dir="ltr">git commit -m "..."</code> فقط تغییر الف ثبت می‌شود، چون <code dir="ltr">git commit</code> محتوای <span dir="ltr">index</span> را ذخیره می‌کند. تغییر ب به‌صورت <span dir="ltr">unstaged</span> باقی می‌ماند و <code dir="ltr">tmp/</code> هم چون <span dir="ltr">untracked</span> است وارد نمی‌شود.</p>

<h3 align="right" dir="rtl">۳. جلوگیری از اضافه شدن <code dir="ltr">tmp/</code></h3>

<p align="right" dir="rtl">خط <code dir="ltr">tmp/</code> را در <code dir="ltr">.gitignore</code> می‌گذارم. اما این فقط روی فایل‌های <span dir="ltr">untracked</span> اثر دارد؛ برای فایلی که قبلاً کامیت شده باید اول آن را از <span dir="ltr">index</span> حذف کرد (فایل روی دیسک می‌ماند):</p>

</div>

```bash
git rm --cached -- draft.txt
```

<div dir="rtl" align="right">

<h2 align="right" dir="rtl">سؤال ۵ — رفع <span dir="ltr">merge conflict</span> در کار گروهی (۴ نمره)</h2>

<p align="right" dir="rtl">من در <code dir="ltr">feature-description</code> و هم‌تیمی‌ام در <code dir="ltr">main</code> یک خط از <code dir="ltr">README.md</code> را متفاوت تغییر داده‌ایم و با <code dir="ltr">git merge feature-description</code> تعارض پیش آمده.</p>

<h3 align="right" dir="rtl">۱. علت</h3>

<p align="right" dir="rtl">در ادغام سه‌طرفه، یک خط در هر دو شاخه نسبت به جد مشترک به دو شکل متفاوت عوض شده و گیت نمی‌تواند تصمیم بگیرد کدام درست است. تغییر بخش‌های جدا از یک فایل تعارض نمی‌سازد؛ فقط تغییر خط‌های یکسان یا مجاور (یا حذف فایلی که طرف دیگر ویرایش کرده).</p>

<h3 align="right" dir="rtl">۲. پیدا کردن تعارض</h3>

<p align="right" dir="rtl"><code dir="ltr">git status</code> فایل را زیر <span dir="ltr">Unmerged paths</span> نشان می‌دهد. داخل فایل:</p>

</div>

```text
<<<<<<< HEAD
This project suggests courses based on student grades.
=======
This project suggests courses using enrollment history.
>>>>>>> feature-description
```

<div dir="rtl" align="right">

<p align="right" dir="rtl">بین <code dir="ltr">&lt;&lt;&lt;&lt;&lt;&lt;&lt; HEAD</code> و <code dir="ltr">=======</code> نسخۀ <code dir="ltr">main</code> است و بین <code dir="ltr">=======</code> و <code dir="ltr">&gt;&gt;&gt;&gt;&gt;&gt;&gt;</code> نسخۀ شاخۀ ادغام‌شونده.</p>

<h3 align="right" dir="rtl">۳. حفظ هر دو تغییر</h3>

<p align="right" dir="rtl">فایل را ویرایش می‌کنم، دو متن را ترکیب و همۀ نشانه‌ها را حذف می‌کنم:</p>

</div>

```text
This project suggests courses based on student grades and enrollment history.
```

```bash
git diff --check
grep -nE '^(<<<<<<<|=======|>>>>>>>)' README.md
git add README.md
git commit
```

<div dir="rtl" align="right">

<p align="right" dir="rtl">برای اطمینان از اینکه هیچ <span dir="ltr">conflict marker</span> نمانده، <code dir="ltr">git diff --check</code> نشانه‌های باقی‌مانده را گزارش می‌کند و خروجی <code dir="ltr">grep</code> باید خالی باشد. بعد <code dir="ltr">git add</code> حل شدن تعارض را اعلام می‌کند و <code dir="ltr">git commit</code> کامیت ادغام را می‌سازد.</p>

<h3 align="right" dir="rtl">۴. لغو <span dir="ltr">merge</span></h3>

<p align="right" dir="rtl"><code dir="ltr">git merge --abort</code> همه‌چیز را به وضعیت پیش از <code dir="ltr">merge</code> برمی‌گرداند (تا وقتی کامیت ادغام ساخته نشده).</p>

<h2 align="right" dir="rtl">سؤال ۶ — انتخاب روش مناسب برای بازگرداندن تغییرات (۴ نمره)</h2>

<h3 align="right" dir="rtl">۱. دور ریختن ویرایش <span dir="ltr">stage</span>‌نشده در <code dir="ltr">notes.txt</code></h3>

</div>

```bash
git restore -- notes.txt
```

<div dir="rtl" align="right">

<p align="right" dir="rtl">فقط <span dir="ltr">working tree</span> همین فایل به نسخۀ <span dir="ltr">staging</span> (که با آخرین کامیت یکی است) برمی‌گردد. <span dir="ltr">staging</span>، تاریخچه و فایل‌های دیگر دست نمی‌خورند. این ویرایش قابل بازیابی نیست.</p>

<h3 align="right" dir="rtl">۲. خارج کردن <code dir="ltr">notes.txt</code> از <span dir="ltr">staging</span> با حفظ ویرایش‌ها</h3>

</div>

```bash
git restore --staged -- notes.txt
```

<div dir="rtl" align="right">

<p align="right" dir="rtl">فقط <span dir="ltr">index</span> به <code dir="ltr">HEAD</code> برمی‌گردد؛ ویرایش‌ها روی دیسک می‌مانند و به‌صورت <span dir="ltr">unstaged</span> دیده می‌شوند.</p>

<h3 align="right" dir="rtl">۳. خنثی کردن کامیت <span dir="ltr">push</span>‌شده روی <code dir="ltr">main</code></h3>

</div>

```bash
git revert --no-edit c0ffee1
git push origin main
```

<div dir="rtl" align="right">

<p align="right" dir="rtl"><code dir="ltr">git revert</code> یک کامیت جدید می‌سازد که اثر <code dir="ltr">c0ffee1</code> را معکوس می‌کند، بدون اینکه تاریخچه را تغییر دهد؛ پس <code dir="ltr">push</code> معمولی کافی است.</p>

<p align="right" dir="rtl"><code dir="ltr">git reset --hard</code> و <code dir="ltr">push --force</code> مناسب نیست، چون کامیتی را که هم‌تیمی‌ها دریافت کرده‌اند از تاریخچه حذف می‌کند، تاریخچۀ آن‌ها را ناسازگار می‌کند و ممکن است کار دیگران را روی شاخۀ دور پاک کند.</p>

<h2 align="right" dir="rtl">استفاده از ابزارهای هوش مصنوعی</h2>

<p align="right" dir="rtl">در نوشتن و ویرایش متن پاسخ‌ها از ابزار هوش مصنوعی کمک گرفته شده است. محتوای همۀ پاسخ‌ها را بررسی کرده‌ام و می‌توانم هر بخش آن را توضیح دهم.</p>

</div>
