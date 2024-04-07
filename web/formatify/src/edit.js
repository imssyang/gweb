import { basicSetup } from "codemirror"
import { EditorView, showPanel, keymap} from "@codemirror/view"
import { EditorState, Compartment, StateField, StateEffect } from "@codemirror/state"
import { defaultKeymap } from "@codemirror/commands"
import { search } from "@codemirror/search"
import { language } from "@codemirror/language"
import { htmlLanguage, html } from "@codemirror/lang-html"
import { javascript } from "@codemirror/lang-javascript"

class HelpPanel {
    constructor(editor) {
        this.toggle = StateEffect.define()
        this.state = StateField.define({
            create: () => false,
            update: (value, tr) => {
                for (let e of tr.effects)
                    if (e.is(this.toggle))
                        value = e.value
                return value
            },
            provide: f => showPanel.from(f, on => on ? this.dom : null)
        })
        this.keymap = [{
            key: "F1",
            run: (view) => {
                view.dispatch({
                    effects: this.toggle.of(!view.state.field(this.state))
                })
                editor.resize()
                return true
            }
        }]
        this.theme = EditorView.baseTheme({
            ".cm-help": {
              padding: "3px 5px",
              backgroundColor: "#fff110",
              fontFamily: "monospace",
            }
        })
    }
    dom(view) {
        let dom = document.createElement("div")
        dom.className = "cm-help"
        dom.textContent = "keymap: https://codemirror.net/6/docs/ref/#commands.standardKeymap"
        return { top: true, dom }
    }
    static create(editor) {
        let help = new HelpPanel(editor)
        return [help.state, keymap.of(help.keymap), help.theme]
    }
}

class StatusPanel {
    constructor(options) {
        this.options = options
    }
    formatBytes(bytes, decimals = 2) {
        if (bytes === 0) return '0';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        const value = parseFloat((bytes / Math.pow(k, i)).toFixed(decimals));
        return `${value} ${sizes[i]}`;
    }
    text(state) {
        const context = state.doc.toString()
        const line = state.doc.lineAt(state.selection.main.head)
        const column = state.selection.main.head - line.from
        const size = this.formatBytes(context.length)
        return `size:${size} pos:${line.number},${column}`
    }
    dom(view) {
        let dom = document.createElement("div")
        dom.classList.add(this.options.class)
        dom.style.color = `#34bc99`
        dom.style.fontSize = '.9em'
        dom.textContent = this.text(view.state)
        return dom
    }
    panel(view) {
        let dom = this.dom(view)
        return {
            dom,
            update: (update) => {
                if (update.docChanged || update.selectionSet)
                    dom.textContent = this.text(update.state)
            }
        }
    }
    static create(options) {
        let status = new StatusPanel(options)
        return showPanel.of((view) => {
            return status.panel(view)
        })
    }
}

class SearchPanel {
    static create() {
        return search({ top: true })
    }
}

class EditorTheme {
    static create(height) {
        return EditorView.theme({
            "&": {
                width: "100%",
            },
            ".cm-scroller": {
                overflow: "auto",
                height: `${height}px`,
            },
            ".cm-content, .cm-gutter": {
                minHeight: "200px",
            },
            ".cm-editor": {
                border: "none",
                fontSize: "14px",
            },
            ".cm-panel.cm-search": {
                "& .cm-textfield": {
                    fontSize: "14px",
                },
                "& .cm-button": {
                    fontSize: "14px",
                    padding: "0 .3em",
                },
                "& label": {
                    "& input": {
                        verticalAlign: "-.2em",
                    }
                }
            }
        }, { dark: false })
    }
}

class EditorLanguage {
    constructor() {
        this.lang = new Compartment
    }
    get default() {
        return this.lang.of(javascript())
    }
    get auto() {
        return EditorState.transactionExtender.of(tr => {
            if (!tr.docChanged) return null
            let docIsHTML = /^\s*</.test(tr.newDoc.sliceString(0, 100))
            let stateIsHTML = tr.startState.facet(language) == htmlLanguage
            if (docIsHTML == stateIsHTML) return null
            return {
                effects: this.lang.reconfigure(docIsHTML ? html() : javascript())
            }
        })
    }
    static create() {
        let language = new EditorLanguage()
        return [language.default, language.auto]
    }
}

class EditorTabSize {
    constructor() {
        this.size = new Compartment()
    }
    static set(size) {
        let ts = EditorTabSize.instance.size
        return ts.reconfigure(EditorState.tabSize.of(size))
    }
    static create() {
        let ts = EditorTabSize.instance.size
        return ts.of(EditorState.tabSize.of(4))
    }
}
EditorTabSize.instance = new EditorTabSize()

class EditorExtensions {
    static create(options) {
        return [
            keymap.of(defaultKeymap),
            EditorView.lineWrapping,
            EditorTheme.create(options.height),
            EditorTabSize.create(),
            ...EditorLanguage.create(),
            ...HelpPanel.create(options.editor),
            SearchPanel.create(),
            StatusPanel.create({
                class: "cm-status",
            })
        ]
    }
}

class CodeEditor {
    constructor(panel, height) {
        this.panel = panel
        this.height = height
        this.dom = (() => {
            const dom = document.createElement("div")
            dom.id = `editor_${panel}`
            return dom
        })()
        this.view = null
    }
    render() {
        this.view = new EditorView({
            parent: document.getElementById(this.dom.id),
            state: EditorState.create({
                doc: '\n'.repeat(0),
                extensions: [
                    basicSetup,
                    EditorExtensions.create({
                        editor: this,
                        height: this.height,
                    }),
                    EditorView.updateListener.of(v => {
                        if (v.heightChanged || v.geometryChanged || v.focusChanged) {
                            this.resize()
                        }
                    })
                ]
            })
        })
    }
    get(id) {
        switch (id) {
            case "html":
                return this.dom.outerHTML
            case "doc":
                return this.view.state.doc.toString().trim()
            default:
                return null
        }
    }
    set(options) {
        let specs = {}
        if (options?.tabSize)
            specs.effects = EditorTabSize.set(options.tabSize)
        if (options.hasOwnProperty('doc'))
            specs.changes = { from: 0, to: this.view.state.doc.length, insert: options.doc }
        this.view.dispatch(specs)
    }
    resize(height) {
        if (height) { this.height = height }
        let help = this.view.dom.querySelector(".cm-help")
        let helpHeight = help ? help.offsetHeight : 0
        let search = this.view.dom.querySelector(".cm-search")
        let searchHeight = search ? search.offsetHeight + 1 : 0
        let status = this.view.dom.querySelector(".cm-status")
        let statusHeight = status ? status.offsetHeight : 0
        let scroller = this.view.dom.querySelector(".cm-scroller")
        let scrollerHeight = this.height - helpHeight - searchHeight - statusHeight
        scroller.style.height = `${scrollerHeight}px`
    }
}

export { CodeEditor }