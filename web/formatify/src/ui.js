import JSON5 from 'json5'
import { w2layout, w2toolbar, w2utils, query } from 'w2ui'
import { CodeEditor } from './edit.js'

class u2toolbar {
    constructor(ui) {
        this.ui = ui
        this.x = new w2toolbar({
            name: 'x2toolbar',
            items: [
                { type: 'menu', id: 'layout',
                    icon: (item) => {
                        return item.get(item.selected)?.icon
                    },
                    selected: '1',
                    items: [
                        { id: '1', text: ' ', icon: 'bi bi-app' },
                        { id: '2', text: ' ', icon: 'bi bi-layout-split' },
                        { id: '3', text: ' ', icon: 'bi bi-layout-three-columns' },
                    ]
                },
                { type: 'break' },
                { type: 'spacer' },
            ],
            onClick: (event) => {
                this.show(event.target)
            }
        })
        this.render()
    }
    get height() {
        return this.x.box.offsetHeight
    }
    render() {
        const dom = document.createElement('div')
        dom.id = this.x.name
        document.body.appendChild(dom)
        this.x.render(`#${dom.id}`)
    }
    show(id) {
        switch (id) {
            case 'layout:1':
                this.ui.layout.show('left', '100%', 1, true)
                this.ui.layout.hide('main')
                this.ui.layout.hide('right')
                break
            case 'layout:2':
                this.ui.layout.show('left', '50%', 2, false)
                this.ui.layout.show('main', '50%', 2, true)
                this.ui.layout.hide('right')
                break
            case 'layout:3':
                this.ui.layout.show('left', '33%', 3, false)
                this.ui.layout.show('main', '34%', 3, false)
                this.ui.layout.show('right', '33%', 3, true)
                break
        }
    }
}

class u2panel {
    constructor(layout, name) {
        this.layout = layout
        this.name = name
        this.id = { left: 1, main: 2, right: 3 }[name].toString()
        this.border = '1px solid gray'
        this.toolbar_meta = {
            style: this.style(true),
            items: [
                { type: 'menu', id: 'index',
                    icon: (item) => {
                        return item.get(this.id)?.icon
                    },
                    selected: this.id,
                    items: [
                        { id: '1', text: ' ', icon: 'bi bi-1-circle', hidden: this.id == '1' },
                        { id: '2', text: ' ', icon: 'bi bi-2-circle', hidden: this.id == '2' },
                        { id: '3', text: ' ', icon: 'bi bi-3-circle', hidden: this.id == '3' },
                    ]
                },
                { type: 'menu', id: 'mode',
                    text: (item) => {
                        return item.get(item.selected)?.text
                    },
                    selected: 'json',
                    items: [
                        { id: 'command', text: 'cmd' },
                        { id: 'json', text: 'json' },
                        { id: 'python', text: 'python' },
                    ]
                },
                { type: 'radio', id: 'contract', group: '1', icon: 'bi bi-chevron-contract' },
                { type: 'radio', id: 'expand', group: '1', icon: 'bi bi-chevron-expand' },
            ],
            onClick: (event) => {
                switch (event.target) {
                    case 'index:1':
                        this.swapDoc('left')
                        break
                    case 'index:2':
                        this.swapDoc('main')
                        break
                    case 'index:3':
                        this.swapDoc('right')
                        break
                    case 'contract':
                    case 'expand':
                        let mode = this.toolbar.get('mode').selected
                        if (mode == 'json') {
                            this.json(event.target == 'expand' ? 4 : 0)
                        } else {
                            let action = event.target
                            let url = `/formatify/${mode}/${action}`
                            this.request(url)
                        }
                        break
                }
            }
        }
        this.editor = null
    }
    get toolbar() {
        return this.layout.x.get(this.name).toolbar
    }
    get height() {
        return this.layout.height - this.layout.ui.toolbar.height - 2
    }
    render() {
        this.editor = new CodeEditor(this.name, this.height)
        this.layout.x.html(this.name, this.editor.get('html'))
        this.editor.render()
    }
    notify(text, error) {
        if (text) {
            w2utils.notify(text, {
                class: 'u2notify',
                timeout: 30_000,
                error: error ? true : false,
            })
        } else {
            query(document.body).find('#w2ui-notify').remove()
        }
    }
    style(toolbar) {
        let v = `border: ${this.border}; border-left: none;`
        if (this.name == 'right')
            v += 'border-right: none;'
        if (toolbar) {
            v += 'border-bottom: none;'
        } else {
            v += 'border-top: none; border-bottom: none;'
        }
        return v
    }
    swapDoc(name) {
        let origin = this.editor
        let target = this.layout.panels[name].editor
        if (origin && target) {
            let odoc = origin.get('doc')
            let tdoc = target.get('doc')
            origin.set({ doc: tdoc })
            target.set({ doc: odoc })
        }
    }
    parseDoc(doc) {
        let obj = null
        try {
            obj = JSON5.parse(doc);
        } catch (error) {
            try {
                obj = Function(`return ${doc}`)()
            } catch (error) {
                this.notify(`${error}`)
            }
        }
        return obj
    }
    json(indent) {
        let doc = this.editor.get('doc')
        let obj = this.parseDoc(doc)
        if (obj) {
            let data = JSON.stringify(obj, null, indent)
            this.editor.set({ doc: data })
            this.notify()
        }
    }
    request(url) {
        let doc = this.editor.get('doc')
        console.log('req', url, JSON.stringify(doc))
        fetch(url, {
            method: 'POST',
            body: JSON.stringify(doc),
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => {
            if (response.ok) {
                return response.json()
            } else {
                throw new Error('HTTP status = ' + response.status)
            }
        })
        .then(data => {
            console.log('rsp', url, data)
        })
        .catch(error => {
            console.error(error)
        })
    }
}

class u2layout {
    constructor(ui) {
        this.ui = ui
        this.immediate = true
        this.panels = {
            main: new u2panel(this, 'main'),
            left: new u2panel(this, 'left'),
            right: new u2panel(this, 'right'),
        }
        this.x = (() => {
            let initPanel = (name) => {
                let panel = this.panels[name]
                return {
                    type: name,
                    size: name == 'main' ? '34%' : '33%',
                    style: panel.style(false),
                    toolbar: panel.toolbar_meta,
                    resizable: true,
                    hidden: name == 'main' ? false : true,
                }
            }
            return new w2layout({
                name: 'x2layout',
                padding: 0,
                panels: [ initPanel('main'), initPanel('left'), initPanel('right') ],
            })
        })()
        this.render()
    }
    get height() {
        return window.innerHeight - this.ui.toolbar.height
    }
    resize() {
        this.x.box.style.height = `${this.height}px`
    }
    render() {
        const dom = document.createElement('div')
        dom.id = this.x.name
        dom.style = `width: 100%; height: ${this.height}px;`
        document.body.appendChild(dom)

        this.x.render(`#${dom.id}`)
        this.x.on('*', (event) => {
            switch (event.type) {
                case 'show':
                    if (['main', 'left', 'right'].includes(event.target)) {
                        let panel = this.panels[event.target]
                        if (!panel.editor) {
                            panel.render()
                        }
                    }
                    break
            }
        })
    }
    set(name, options) {
        let panel = this.panels[name]
        if (options.hasOwnProperty('last')) {
            let borderRight = options.last ? 'none' : panel.border
            let resizable = options.last ? false : true
            this.x.get(name).resizable = resizable
            this.x.get(name).toolbar.box.style.borderRight = borderRight
            this.x.el(name).style.borderRight = borderRight

            if (name == 'left') {
                let resizerId = `#layout_${this.x.name}_resizer_${name}`
                let resizer = query(this.x.box).find(resizerId)
                if (resizable) {
                    resizer.show()
                } else {
                    resizer.hide()
                }
            }
        }
        if (options?.layout) {
            let index = this.x.get(name).toolbar.get('index')
            Object.values([1, 2, 3]).forEach(id => {
                let sameId = panel.id == id.toString()
                index.get(id).hidden = sameId || id > options.layout
            })
        }
    }
    show(name, size, layout, last) {
        this.x.sizeTo(name, size, this.immediate)
        this.x.show(name, this.immediate)
        this.set(name, {
            layout: layout,
            last: last,
        })
    }
    hide(name) {
        this.x.hide(name, this.immediate)
    }
}

class UI {
    constructor() {
        this.toolbar = new u2toolbar(this)
        this.layout = new u2layout(this)
        this.toolbar.show('layout:1')
        this.resize()
    }
    resize() {
        window.addEventListener('resize', () => {
            this.layout.resize()
            Object.values(this.layout.panels).forEach(panel => {
                let editor = panel.editor
                if (editor instanceof CodeEditor) {
                    editor.resize(panel.height)
                }
            })
        })
    }
}

export { UI }