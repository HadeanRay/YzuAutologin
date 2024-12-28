import './style.css';

import {SaveValue, ReadData, Loginyzu, EnableAutoStart, DisableAutoStart} from '../wailsjs/go/main/App';
import { Hide } from '../wailsjs/runtime/runtime';
import 'sober';

let webindex = document.getElementById("webindex");
let countindex = document.getElementById("countindex");
let passwordindex = document.getElementById("passwordindex");
let operatorindex = document.getElementById("operatorindex");
let autostartindex = document.getElementById("autostartindex");
let operatorindex_items = document.querySelectorAll('s-segmented-button-item');
let testconnectindex = document.getElementById("testconnectindex");

// 设置一个定时器变量
let typingTimer;
let doneTypingInterval = 500; // 时间间隔（毫秒）

// 用户停止输入后的处理函数
function doneTyping() {

    try {
        SaveValue({
            [webindex.id]: webindex.value,
            [countindex.id]: countindex.value,
            [passwordindex.id]: passwordindex.value,
            [operatorindex.id]: operatorindex.value,
            [autostartindex.id]: autostartindex.checked.toString()
        })
    } catch (err) {
        console.error(err); 
    }
}

// 在用户输入时清除定时器
webindex.addEventListener('input', () => {
    clearTimeout(typingTimer);
    typingTimer = setTimeout(doneTyping, doneTypingInterval);
});

countindex.addEventListener('input', () => {
    clearTimeout(typingTimer);
    typingTimer = setTimeout(doneTyping, doneTypingInterval);
});

passwordindex.addEventListener('input', () => {
    clearTimeout(typingTimer);
    typingTimer = setTimeout(doneTyping, doneTypingInterval);
});

operatorindex.addEventListener('click', () => {
    clearTimeout(typingTimer);
    typingTimer = setTimeout(doneTyping, doneTypingInterval);
});

autostartindex.addEventListener('click', () => {
    clearTimeout(typingTimer);
    typingTimer = setTimeout(doneTyping, doneTypingInterval);
    try {
        if (autostartindex.checked) {
            EnableAutoStart();
        } else {
            DisableAutoStart();
        }
    } catch (err) {
        console.error('Failed to enable auto start:', err);
    }
}); 

testconnectindex.addEventListener('click', async () => {
    try {
        await Loginyzu();
        showSnackbar("测试连接已触发");
    } catch (err) {
        console.error(err);
        showSnackbar(err.toString());
    }
});

ReadData()
    .then((data) => {
        webindex.value = data.webindex;
        countindex.value = data.countindex;
        passwordindex.value = data.passwordindex;
        operatorindex.value = data.operatorindex;
        autostartindex.checked = data.autostartindex == "true";
        
        if (data.autostartindex == "true") {
            (async () => {
                try {
                    await Loginyzu();
                    showSnackbar("测试连接已触发");
                } catch (err) {
                    console.error(err);
                    showSnackbar(err.toString());
                }
            })();
        }

        setTimeout(() => {
            operatorindex_items.forEach((item) => {
                forceRedraw(item);
            });
        }, 1000);
        
    })
    .catch((err) => {
        console.error("Error reading data:", err);
    }
);

let exitelement = document.getElementById("exit");

exitelement.addEventListener('click', () => {
    try {
        Hide();
    } catch (err) {
        console.error(err); 
    }
});

// 强制重绘 s-segmented-button 元素
function forceRedraw(element) {
    element.style.display = 'none';
    element.offsetHeight; // 触发重绘
    element.style.display = '';
};

// 显示 Snackbar 消息通知
function showSnackbar(message) {
    const snackbar = document.createElement('s-snackbar');
    const htmlContent = `
        <s-button slot="trigger" class="s-button--text" style="background-color: transparent"></s-button>
        ${message}
    `;
    snackbar.innerHTML = htmlContent; 
    
    const sPage = document.querySelector('s-page');
    sPage.appendChild(snackbar);

    snackbar.querySelector('s-button').click();
    setTimeout(() => {
        snackbar.remove();
    }, 5000); 
}
