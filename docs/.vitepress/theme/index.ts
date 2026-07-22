import DefaultTheme from 'vitepress/theme'
import type { Theme } from 'vitepress'
import DocFigure from './components/DocFigure.vue'
import HomeHero from './components/HomeHero.vue'
import LearningPath from './components/LearningPath.vue'
import MermaidDiagram from './components/MermaidDiagram.vue'
import PracticeBlock from './components/PracticeBlock.vue'
import TechGrid from './components/TechGrid.vue'
import './styles.css'

export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    app.component('DocFigure', DocFigure)
    app.component('HomeHero', HomeHero)
    app.component('LearningPath', LearningPath)
    app.component('MermaidDiagram', MermaidDiagram)
    app.component('PracticeBlock', PracticeBlock)
    app.component('TechGrid', TechGrid)
  }
} satisfies Theme
