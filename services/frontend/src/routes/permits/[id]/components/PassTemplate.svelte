<script lang="ts">
    /* Використовуємо Svelte 5 Runes */
    let { permit = {} } = $props();

    // Дані з картинки як фолбек
    const displayData = $derived({
        date: permit.entry_time ? new Date(permit.entry_time).toLocaleDateString('uk-UA') : "",
        time: permit.entry_time ? new Date(permit.entry_time).toLocaleTimeString('uk-UA', {hour: '2-digit', minute:'2-digit'}) : "",
        id: permit.ID ?? "",
        code: permit.code ?? "",
        plate1: permit.plate_front ?? "",
        plate2: permit.plate_back ?? "",
        post: permit.customs_post?.name ?? ""
    });
</script>

<div class="a6-wrapper">
    <div class="a6-page">
        <div class="outer-border">
            <div class="inner-content">
                <div class="header-grid">
                    <div class="meta">
                        <div class="main-bold">{displayData.date}</div>
                        <div class="main-bold">{displayData.time}</div>
                        <div class="id-num">{displayData.id}</div>
                    </div>

                    <div class="title-area">
                        <h1 class="pass-num">ПЕРЕПУСТКА № {displayData.code}</h1>
                        <p class="desc">Видана транспортному засобу з реєстраційним номером</p>
                        <div class="plates-row">
                            <span>{displayData.plate1}</span>
                            <span>{displayData.plate2}</span>
                        </div>
                    </div>

                    <div class="barcode-area">
                        <div class="barcode-mock">
                            {#each Array(35) as _, i}
                                <div class="bar" style="width: {i % 4 === 0 ? '3.5px' : '1.5px'}"></div>
                            {/each}
                        </div>
                    </div>
                </div>

                <div class="line-separator"></div>

                <div class="footer-grid">
                    <div class="side-text">
                        Дозволено в'їзд у зону митного контролю
                    </div>

                    <div class="stamp-box">
                        <div class="stamp-label">МП</div>
                        <div class="stamp-value">{displayData.post}</div>
                    </div>

                    <div class="side-text">
                        Дозволено виїзд за межі зони митного контролю
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    /* ПРИХОВУЄМО НА ЕКРАНІ */
    .a6-wrapper {
        display: none;
    }

    /* ДРУК */
    @media print {
        @page {
            size: A6 landscape;
            margin: 0;
        }

        /* Повне очищення сторінки */
        :global(html), :global(body) {
            margin: 0 !important;
            padding: 0 !important;
            height: 105mm !important;
            width: 148mm !important;
            overflow: hidden !important;
            background: white !important;
        }

        /* Приховуємо все по замовчуванню */
        :global(body *) {
            visibility: hidden !important;
            height: 0 !important;
            margin: 0 !important;
            padding: 0 !important;
            border: none !important;
        }

        /* Показуємо лише перепустку та її вміст */
        .a6-wrapper, .a6-wrapper * {
            visibility: visible !important;
            height: auto !important;
        }

        .a6-wrapper {
            display: block !important;
            position: fixed !important;
            top: 0 !important;
            left: 0 !important;
            width: 148mm !important;
            height: 105mm !important;
            z-index: 999999 !important;
            background: white !important;
            padding: 0 !important;
            margin: 0 !important;
            border: none !important;
        }

        .a6-page {
            width: 148mm !important;
            height: 105mm !important;
            padding: 4mm !important;
            box-sizing: border-box !important;
            background: white !important;
            font-family: Arial, Helvetica, sans-serif !important;
            color: black !important;
            display: flex !important;
            flex-direction: column !important;
        }

        .outer-border {
            border: 1px solid black !important;
            width: 100% !important;
            height: 100% !important;
            padding: 2mm !important;
            box-sizing: border-box !important;
        }

        .inner-content {
            border: 1px solid black !important;
            width: 100% !important;
            height: 100% !important;
            padding: 3mm !important;
            box-sizing: border-box !important;
            display: flex !important;
            flex-direction: column !important;
        }

        .header-grid {
            display: grid !important;
            grid-template-columns: 1fr 2.5fr 1fr !important;
            gap: 2mm !important;
        }

        .main-bold {
            font-size: 15pt !important;
            font-weight: 950 !important;
            line-height: 1.1 !important;
        }

        .id-num {
            font-size: 13pt !important;
            margin-top: 5mm !important;
            font-weight: bold !important;
        }

        .title-area {
            text-align: center !important;
        }

        .pass-num {
            font-size: 18pt !important;
            font-weight: 950 !important;
            margin: 0 !important;
            white-space: nowrap !important;
        }

        .desc {
            font-size: 10pt !important;
            font-weight: bold !important;
            margin: 2mm 0 !important;
        }

        .plates-row {
            display: flex !important;
            justify-content: center !important;
            gap: 8mm !important;
            font-size: 18pt !important;
            font-weight: 900 !important;
        }

        .barcode-area {
            display: flex !important;
            justify-content: flex-end !important;
        }

        .barcode-mock {
            display: flex !important;
            height: 12mm !important;
            gap: 1.5px !important;
            align-items: stretch !important;
        }

        .bar {
            background: black !important;
            height: 100% !important;
        }

        .line-separator {
            border-top: 1.2px solid black !important;
            margin: 4mm 0 !important;
        }

        .footer-grid {
            display: grid !important;
            grid-template-columns: 1fr auto 1fr !important;
            align-items: center !important;
            margin-top: auto !important;
            padding-bottom: 3mm !important;
        }

        .side-text {
            font-size: 11pt !important;
            font-weight: 900 !important;
            text-align: center !important;
            line-height: 1.2 !important;
            padding: 0 3mm !important;
        }

        .stamp-box {
            border: 1.2px solid black !important;
            padding: 1.5mm 4mm !important;
            text-align: center !important;
            min-width: 30mm !important;
            height: auto !important;
        }

        .stamp-label { font-size: 8pt !important; font-weight: bold !important; }
        .stamp-value { font-size: 11pt !important; font-weight: 950 !important; text-transform: uppercase !important; }
    }
</style>