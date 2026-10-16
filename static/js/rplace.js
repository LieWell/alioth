(function () {
  'use strict';

  const W = 32, H = 16;
  const WHITE = '#FAFAFA', BLACK = '#404040';

  // 内存中的像素状态：{color, rotation}
  const pixels = Array.from({ length: W * H }, () => ({ color: WHITE, rotation: 0 }));

  // 渲染 32x16 网格
  const canvas = document.getElementById('canvas');
  canvas.style.gridTemplateColumns = `repeat(${W}, 20px)`;
  for (let i = 0; i < pixels.length; i++) {
    const el = document.createElement('div');
    el.className = 'pixel';
    el.style.backgroundColor = pixels[i].color;
    el.addEventListener('click', () => updatePixel(i));
    canvas.appendChild(el);
  }

  // WebSocket：根据当前协议自动选择 ws/wss
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
  const socket = new WebSocket(`${proto}//${location.host}/rplace`);

  socket.onopen = () => console.log('[rplace] ws connected');
  socket.onerror = (e) => console.error('[rplace] ws error', e);

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (Array.isArray(data)) {
      // 建立连接后服务端推送的整张画布
      data.forEach((row, y) => row.forEach((c, x) => applyPixel(x, y, c)));
    } else if (data && typeof data.x === 'number') {
      // 增量更新
      applyPixel(data.x, data.y, data.color);
    }
  };

  function applyPixel(x, y, intColor) {
    const idx = y * W + x;
    if (!pixels[idx]) return;
    pixels[idx].color = intColor === 1 ? WHITE : BLACK;
    const el = canvas.children[idx];
    if (el) el.style.backgroundColor = pixels[idx].color;
  }

  function updatePixel(idx) {
    const x = idx % W;
    const y = Math.floor(idx / W);
    const el = canvas.children[idx];
    const newColor = pixels[idx].color === BLACK ? WHITE : BLACK;
    const rotation = pixels[idx].rotation === 90 ? 0 : 90;

    pixels[idx].color = newColor;
    pixels[idx].rotation = rotation;
    el.style.backgroundColor = newColor;
    el.style.transform = `rotate(${rotation}deg)`;

    socket.send(JSON.stringify({ x: x, y: y, color: newColor === WHITE ? 1 : 0 }));
  }
})();
