--------------------------------------------------------------------------------
{-# LANGUAGE OverloadedStrings #-}
import           Data.Monoid (mappend)
import           Hakyll
import           Text.Pandoc.Options
import           Text.Pandoc.Highlighting (styleToCss, haddock)
--------------------------------------------------------------------------------
main :: IO ()
main = hakyll $ do
    match (
        "images/*"  .||. 
        "fonts/*"   .||. 
        "docs/*"    .||.
        "favicon.ico") $ do
        route   idRoute
        compile copyFileCompiler

    match "css/*" $ do
        route   idRoute
        compile compressCssCompiler

    create ["css/syntax.css"] $ do
        route idRoute
        compile $ makeItem $ styleToCss haddock

    match "posts/*" $ do
        route $ setExtension "html"
        compile $ pandocCompilerWith defaultHakyllReaderOptions
            defaultHakyllWriterOptions { writerHighlightStyle = Just haddock }
            >>= loadAndApplyTemplate "templates/post.html"    postCtx
            >>= loadAndApplyTemplate "templates/default.html" postCtx
            >>= relativizeUrls

    create ["archive.html"] $ do
        route idRoute
        compile $ do
            posts <- recentFirst =<< loadAll "posts/*"
            let archiveCtx =
                    listField "posts" postCtx (return posts) `mappend`
                    constField "title" "blog archives"       `mappend`
                    constField "isBlog" "true"               `mappend`
                    defaultContext

            makeItem ""
                >>= loadAndApplyTemplate "templates/archive.html" archiveCtx
                >>= loadAndApplyTemplate "templates/default.html" archiveCtx
                >>= relativizeUrls

    match "index.html" $ do
        route idRoute
        compile $ do
            posts <- recentFirst =<< loadAll "posts/*"
            let posts' = take 4 posts
                indexCtx =
                    listField "posts" postCtx (return posts') `mappend`
                    constField "title" "jackie's site!"       `mappend`
                    constField "isHome" "true"                `mappend`
                    defaultContext

            getResourceBody
                >>= applyAsTemplate indexCtx
                >>= loadAndApplyTemplate "templates/default.html" indexCtx
                >>= relativizeUrls

    match "templates/*" $ compile templateBodyCompiler
--------------------------------------------------------------------------------
postCtx :: Context String
postCtx =
    dateField "date" "%d %b %Y" `mappend`
    constField "isBlog" "true"  `mappend`
    defaultContext
