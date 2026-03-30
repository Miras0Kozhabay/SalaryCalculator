const API_BASE = '';
    let currentMode = 'gross';
    let currentPage = 0;
    const PAGE_SIZE = 10;
    let totalHistory = 0;

    // Theme
    let isDark = localStorage.getItem('theme') !== 'light';
    function applyTheme() {
      document.documentElement.classList.toggle('light', !isDark);
      document.getElementById('themeToggle').textContent = isDark ? '🌙' : '☀️';
    }
    function toggleTheme() {
      isDark = !isDark;
      localStorage.setItem('theme', isDark ? 'dark' : 'light');
      applyTheme();
    }
    applyTheme();

    function setMode(mode) {
      currentMode = mode;
      document.querySelectorAll('.mode-btn').forEach(b => {
        b.classList.toggle('active', b.dataset.mode === mode);
      });
      document.getElementById('modeHint').textContent =
        mode === 'gross'
          ? 'Рассчитать NET из GROSS (до вычета налогов)'
          : 'Рассчитать GROSS из NET (сумма на руки)';
    }

    function fmt(n) {
      if (n == null || isNaN(n)) return '—';
      return new Intl.NumberFormat('ru-KZ', { maximumFractionDigits: 2 }).format(n);
    }

    function showError(msg) {
      const el = document.getElementById('errorBanner');
      el.textContent = '⚠ ' + msg;
      el.classList.add('show');
    }

    function hideError() {
      document.getElementById('errorBanner').classList.remove('show');
    }

    function setLoading(v) {
      const btn = document.getElementById('calcBtn');
      document.getElementById('btnText').style.display = v ? 'none' : '';
      document.getElementById('btnSpinner').style.display = v ? '' : 'none';
      btn.disabled = v;
    }

    async function calculate() {
      hideError();
      const raw = document.getElementById('salaryInput').value.trim();
      const salary = parseFloat(raw);
      const input = document.getElementById('salaryInput');

      // Validate input
      if (!raw || isNaN(salary) || salary <= 0) {
        input.classList.add('error');
        showError('Введите корректную сумму (больше 0)');
        return;
      }
      input.classList.remove('error');

      setLoading(true);
      
      // Create abort controller for timeout
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 10000); // 10 second timeout

      try {
        const res = await fetch(`${API_BASE}/api/calculate`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ salary, mode: currentMode }),
          signal: controller.signal
        });

        if (!res.ok) {
          let errorMsg;
          try {
            const err = await res.json();
            errorMsg = err.error || err.message || `HTTP ${res.status}`;
          } catch {
            errorMsg = `Server error: HTTP ${res.status}`;
          }
          throw new Error(errorMsg);
        }

        const data = await res.json();
        renderResults(data);
        input.value = ''; // Clear input after successful calculation
        loadHistory(0);
      } catch (e) {
        // Handle abort (timeout)
        if (e.name === 'AbortError') {
          showError('Request timeout - server not responding (10s)');
        } else if (e.message === 'Failed to fetch') {
          showError('Network error - please check your connection');
        } else {
          showError(e.message);
        }
      } finally {
        clearTimeout(timeoutId);
        setLoading(false);
      }
    }

    function renderResults(d) {
      const gross = d.gross_salary ?? d.gross ?? 0;
      const net = d.net_salary ?? d.net ?? 0;
      const opv = d.opv ?? 0;
      const vosms = d.vosms ?? 0;
      const ipn = d.ipn ?? 0;
      const so = d.so ?? 0;
      const oosms = d.oosms ?? 0;
      const sn = d.sn ?? 0;
      const employerTotal = d.employer_total ?? (gross + so + oosms + sn);

      document.getElementById('resultsBody').innerHTML = `
        <div class="fade-in">
          <div class="result-net">
            <div class="result-net-label">💼 Зарплата на руки (NET)</div>
            <div class="result-net-value">${fmt(net)}<span>₸</span></div>
          </div>

          <div class="section-label">Вычеты из зарплаты сотрудника</div>
          <div class="tax-grid">
            <div class="tax-row">
              <span class="tax-label"><span class="abbr">GROSS</span> До вычетов</span>
              <span class="tax-value">${fmt(gross)} ₸</span>
            </div>
            <div class="tax-row">
              <span class="tax-label"><span class="abbr">ОПВ</span> Пенсионный взнос</span>
              <span class="tax-value neg">−${fmt(opv)} ₸</span>
            </div>
            <div class="tax-row">
              <span class="tax-label"><span class="abbr">ВОСМС</span> Мед. страхование</span>
              <span class="tax-value neg">−${fmt(vosms)} ₸</span>
            </div>
            <div class="tax-row">
              <span class="tax-label"><span class="abbr">ИПН</span> Подоходный налог</span>
              <span class="tax-value neg">−${fmt(ipn)} ₸</span>
            </div>
          </div>

          <div class="section-label">Отчисления работодателя</div>
          <div class="tax-grid">
            <div class="tax-row">
              <span class="tax-label"><span class="abbr">СО</span> Социальные отчисления</span>
              <span class="tax-value pos">${fmt(so)} ₸</span>
            </div>
            <div class="tax-row">
              <span class="tax-label"><span class="abbr">ООСМС</span> Взнос работодателя ОСМС</span>
              <span class="tax-value pos">${fmt(oosms)} ₸</span>
            </div>
            <div class="tax-row">
              <span class="tax-label"><span class="abbr">СН</span> Социальный налог</span>
              <span class="tax-value pos">${fmt(sn)} ₸</span>
            </div>
          </div>

          <div class="employer-total">
            <span class="employer-total-label">Итого расходы работодателя</span>
            <span class="employer-total-value">${fmt(employerTotal)} ₸</span>
          </div>
        </div>
      `;
    }

    async function loadHistory(page) {
      if (page !== undefined) currentPage = page;
      const offset = currentPage * PAGE_SIZE;
      try {
        const res = await fetch(`${API_BASE}/api/history?limit=${PAGE_SIZE}&offset=${offset}`);
        if (!res.ok) return;
        const data = await res.json();
        renderHistory(data);
      } catch (e) {
        // silent fail for history
      }
    }

    function renderHistory(data) {
      const items = Array.isArray(data)
        ? data
        : (data.items || data.calculations || data.data || []);
      const total = data.total ?? items.length;
      totalHistory = total;

      const el = document.getElementById('historyBody');

      if (!items.length) {
        el.innerHTML = `<div class="state-empty" style="min-height:140px;"><div class="icon">☰</div><p>История пуста</p></div>`;
        document.getElementById('pagination').style.display = 'none';
        return;
      }

      el.innerHTML = `
        <div style="overflow-x:auto;">
          <table class="history-table">
            <thead>
              <tr>
                <th>#</th>
                <th>Тип</th>
                <th>GROSS</th>
                <th>NET</th>
                <th>ОПВ</th>
                <th>ИПН</th>
                <th>👔 Работодатель</th>
                <th>Дата</th>
              </tr>
            </thead>
            <tbody>
              ${items.map((row, i) => {
                const gross = row.gross_salary ?? row.gross ?? 0;
                const net = row.net_salary ?? row.net ?? 0;
                const opv = row.opv ?? 0;
                const ipn = row.ipn ?? 0;
                const employerTotal = row.employer_total ?? 0;
                const mode = row.mode || 'gross';
                return `
                  <tr>
                    <td class="mono" style="color:var(--text-muted)">${currentPage * PAGE_SIZE + i + 1}</td>
                    <td><span class="badge badge-${mode}">${mode.toUpperCase()}</span></td>
                    <td class="mono">${fmt(gross)} ₸</td>
                    <td class="mono" style="color:var(--accent-blue);font-weight:700">${fmt(net)} ₸</td>
                    <td class="mono" style="color:#f87171">−${fmt(opv)} ₸</td>
                    <td class="mono" style="color:#f87171">−${fmt(ipn)} ₸</td>
                    <td class="mono" style="color:var(--accent-green);font-weight:600">${fmt(employerTotal)} ₸</td>
                    <td class="mono" style="color:var(--text-muted)">${formatDate(row.created_at)}</td>
                  </tr>
                `;
              }).join('')}
            </tbody>
          </table>
        </div>
      `;

      const pag = document.getElementById('pagination');
      const hasPag = total > PAGE_SIZE || currentPage > 0;
      pag.style.display = hasPag ? 'flex' : 'none';

      if (hasPag) {
        const from = currentPage * PAGE_SIZE + 1;
        const to = Math.min((currentPage + 1) * PAGE_SIZE, total);
        document.getElementById('paginationInfo').textContent = `${from}–${to} из ${total}`;
        document.getElementById('prevBtn').disabled = currentPage === 0;
        document.getElementById('nextBtn').disabled = (currentPage + 1) * PAGE_SIZE >= total;
      }
    }

    function changePage(dir) {
      loadHistory(currentPage + dir);
    }

    function formatDate(raw) {
      if (!raw) return '—';
      try {
        const d = new Date(raw);
        return d.toLocaleDateString('ru-KZ', { day: '2-digit', month: '2-digit', year: '2-digit' })
          + ' ' + d.toLocaleTimeString('ru-KZ', { hour: '2-digit', minute: '2-digit' });
      } catch { return String(raw); }
    }

    document.getElementById('salaryInput').addEventListener('keydown', e => {
      if (e.key === 'Enter') calculate();
    });

    document.getElementById('salaryInput').addEventListener('input', () => {
      document.getElementById('salaryInput').classList.remove('error');
      hideError();
    });

    // Load history on startup
    loadHistory(0);