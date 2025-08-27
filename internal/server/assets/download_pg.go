package assets

var DownloadPg = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>文件下载</title>
    <!-- 引入 Tailwind CSS -->
    <script src="/tailwindcss"></script>
    <script>
        tailwind.config = {
            theme: {
                extend: {
                    colors: {
                        primary: '#3b82f6',
                        danger: '#ef4444',
                    },
                }
            }
        }
    </script>
</head>
<body class="bg-gray-50 min-h-screen p-4 sm:p-6">
<div class="max-w-2xl mx-auto mt-10">
    <!-- 文件卡片 -->
    <div class="bg-white rounded-lg border border-gray-100 shadow-sm p-5 sm:p-6">
        <!-- 文件图标与名称 -->
        <div class="flex items-center gap-4 mb-5">
            <div>
                <h1 class="text-xl sm:text-2xl font-semibold text-gray-800">产品需求文档_v2.1.pdf</h1>
            </div>
        </div>

        <!-- 核心文件信息 -->
        <div class="grid grid-cols-2 gap-4 mb-6 text-sm">
            <div class="flex items-center gap-2 text-gray-600">
                <i class="fa fa-hdd-o text-gray-400"></i>
                <span>大小: 2.4 MB</span>
            </div>
            <div class="flex items-center gap-2 text-gray-600">
                <i class="fa fa-file-text-o text-gray-400"></i>
                <span>格式: PDF</span>
            </div>
        </div>

        <!-- 密码输入区域 -->
        <div class="mb-6">
            <p class="text-sm text-gray-600 mb-2 flex items-center">
                <i class="fa fa-lock text-primary mr-1"></i>
                该文件受密码保护，请输入下载密码
            </p>
            <div class="relative">
                <input
                        type="password"
                        id="downloadPassword"
                        placeholder="请输入下载密码"
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition"
                >
            </div>
            <p id="passwordError" class="text-danger text-xs mt-1 hidden">密码错误，请重新输入</p>
        </div>

        <!-- 下载按钮 -->
        <div class="pt-4 border-t border-gray-100">
            <button id="downloadBtn" class="w-full sm:w-auto bg-primary hover:bg-primary/90 text-white py-3 px-6 rounded-lg flex items-center justify-center gap-2 transition-colors duration-200 font-medium">
                <i class="fa fa-download"></i>
                <span>下载文件</span>
            </button>
        </div>
    </div>

    <!-- 底部提示 -->
    <p class="text-center text-gray-500 text-xs mt-4">
        TFLanHttpDesktop v0.1
    </p>
</div>

<script>
    // 实际应用中应从后端获取正确密码（此处仅为示例）
    const correctPassword = "123456"; // 替换为实际密码验证逻辑

    const downloadBtn = document.getElementById('downloadBtn');
    const passwordInput = document.getElementById('downloadPassword');
    const passwordError = document.getElementById('passwordError');

    downloadBtn.addEventListener('click', function(e) {
        e.preventDefault();
        const password = passwordInput.value.trim();

        // 验证密码
        if (password === '') {
            showError('请输入下载密码');
            return;
        }

        if (password !== correctPassword) {
            showError('密码错误，请重新输入');
            return;
        }

        // 密码正确，执行下载流程
        startDownload();
    });

    // 显示错误提示
    function showError(message) {
        passwordError.textContent = message;
        passwordError.classList.remove('hidden');
        passwordInput.classList.add('border-danger');

        // 3秒后自动隐藏错误提示
        setTimeout(() => {
            passwordError.classList.add('hidden');
            passwordInput.classList.remove('border-danger');
        }, 3000);
    }

    // 模拟下载过程
    function startDownload() {
        const originalText = downloadBtn.innerHTML;
        downloadBtn.disabled = true;
        downloadBtn.innerHTML = '<i class="fa fa-spinner fa-spin"></i><span>下载中...</span>';
        downloadBtn.classList.add('opacity-80', 'cursor-wait');

        // 模拟下载延迟
        setTimeout(() => {
            downloadBtn.innerHTML = '<i class="fa fa-check"></i><span>下载完成</span>';
            downloadBtn.classList.remove('bg-primary', 'opacity-80', 'cursor-wait');
            downloadBtn.classList.add('bg-green-500');

            // 实际应用中这里应该是真实的文件下载链接
            // window.location.href = 'path/to/file.pdf';
        }, 1500);
    }

    // 输入密码时清除错误状态
    passwordInput.addEventListener('input', function() {
        passwordError.classList.add('hidden');
        passwordInput.classList.remove('border-danger');
    });
</script>
</body>
</html>

`
