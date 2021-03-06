// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

export type Theme = {
    mainBg: string
    mainFg: string
    buttonBg: string
    buttonFg: string
    sidebarBg: string
    sidebarFg: string
    sidebarWhiteLogo: string
}

export const defaultTheme = {
    mainBg: '255, 255, 255',
    mainFg: '55, 53, 47',
    buttonBg: '22, 109, 224',
    buttonFg: '255, 255, 255',
    sidebarBg: '20, 93, 191',
    sidebarFg: '255, 255, 255',
    sidebarWhiteLogo: 'true',
}

export const darkTheme = {
    mainBg: '55, 53, 47',
    mainFg: '200, 200, 200',
    buttonBg: '80, 170, 221',
    buttonFg: '255, 255, 255',
    sidebarBg: '75, 73, 67',
    sidebarFg: '255, 255, 255',
    sidebarWhiteLogo: 'true',
}

export const lightTheme = {
    mainBg: '255, 255, 255',
    mainFg: '55, 53, 47',
    buttonBg: '80, 170, 221',
    buttonFg: '255, 255, 255',
    sidebarBg: '247, 246, 243',
    sidebarFg: '55, 53, 47',
    sidebarWhiteLogo: 'false',
}

export function setTheme(theme: Theme): void {
    document.documentElement.style.setProperty('--main-bg', theme.mainBg)
    document.documentElement.style.setProperty('--main-fg', theme.mainFg)
    document.documentElement.style.setProperty('--body-color', theme.mainFg)
    document.documentElement.style.setProperty('--button-bg', theme.buttonBg)
    document.documentElement.style.setProperty('--button-fg', theme.buttonFg)
    document.documentElement.style.setProperty('--sidebar-bg', theme.sidebarBg)
    document.documentElement.style.setProperty('--sidebar-fg', theme.sidebarFg)
    document.documentElement.style.setProperty('--sidebar-white-logo', theme.sidebarWhiteLogo)
    localStorage.setItem('theme', JSON.stringify(theme))
}

export function loadTheme(): Theme {
    const themeStr = localStorage.getItem('theme')
    if (themeStr) {
        try {
            const theme = JSON.parse(themeStr)
            setTheme(theme)
            return theme
        } catch (e) {
            setTheme(defaultTheme)
            return defaultTheme
        }
    } else {
        setTheme(defaultTheme)
        return defaultTheme
    }
}
