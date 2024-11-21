"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle,CardDescription,CardFooter } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Clipboard, LinkIcon, ExternalLink } from 'lucide-react'
import axios from "axios"
import { useToast } from "@/hooks/use-toast"

function isValidURL(url: string) {
  const pattern = /^https:\/\/.+/;
  return pattern.test(url);
}

export default function URLShortener() {
  const { toast } = useToast()

  const [longUrl, setLongUrl] = useState("")
  const [shortUrl, setShortUrl] = useState("")
  const[conveting,setConverting]=useState(false)

  const [recentUrls, setRecentUrls] = useState([
    { long: "https://example.com/very/long/url/that/needs/shortening", short: "https://short.url/abc123" },
    { long: "https://another-example.com/with/a/long/url", short: "https://short.url/def456" },
  ])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      setConverting(true)
      console.log("checking")
      toast({
        title: "Converting URL"
      })
      if (!isValidURL(longUrl)) {
        toast({
          title: "Invalid URL",
          description: "Please enter a valid URL",
          variant : "destructive"
        })
        return
      }
      const res = await axios.post("https://short.paramveer.in", {LongUrl: longUrl })
      if(res.data==="Rate limit exceeded"){
        toast({
          title: "Rate Limit Exceeded",
          description: "Please try after 24hr",
          variant : "destructive"
        })
        return
      }
      setShortUrl("https://short.paramveer.in/"+res.data.ShortUrl)
      toast({
        title: "Url Converted"
      })
      console.log(res)
    } catch (error) {
      toast({
        title: "Error",
        description: "An error occurred while converting the URL",
        variant : "destructive"
      })
    } finally{
      setConverting(false)
    }
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-200 to-blue-50 dark:from-gray-800 dark:to-gray-700">
      <header className="bg-blue-100 dark:bg-gray-900 shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <div className="flex items-center">
              <LinkIcon className="h-8 w-8 text-blue-500" />
              <span className="ml-2 text-2xl font-bold text-gray-900 dark:text-white">URL Shortener</span>
            </div>
            <nav>
              <ul className="flex space-x-4">
                <li><a href="#" className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white">Home</a></li>
                <li><a href="#" className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white">About</a></li>
                <li><a href="#" className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white">Contact</a></li>
              </ul>
            </nav>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <Card className="w-full max-w-2xl mx-auto bg-white">
          <CardHeader>
            <CardTitle>Shorten Your URL</CardTitle>
            <CardDescription>Enter a long URL to get a short, shareable link.</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="longUrl">Long URL</Label>
                <Input
                  id="longUrl"
                  placeholder="https://example.com/your/long/url"
                  value={longUrl}
                  onChange={(e) => setLongUrl(e.target.value)}
                  required
                />
              </div>
              <Button type="submit" className="w-full bg-black text-white hover:text-black" disabled={conveting} >Shorten URL</Button>
            </form>
          </CardContent>
          {shortUrl && (
            <CardFooter className="flex flex-col items-start space-y-2">
              <Label htmlFor="shortUrl">Your shortened URL:</Label>
              <div className="flex w-full items-center space-x-2">
                <Input id="shortUrl" value={shortUrl} readOnly />
                <Button size="icon" onClick={() => copyToClipboard(shortUrl)}>
                  <Clipboard className="h-4 w-4" />
                </Button>
              </div>
            </CardFooter>
          )}
        </Card>

        <Tabs className="w-full max-w-2xl mx-auto mt-8">
          <TabsList className="grid w-full grid-cols-1 bg-slate-400">
            <TabsTrigger value="recent" className="rounded-lg">Recent URLs</TabsTrigger>
          </TabsList>
          <TabsContent value="recent" className="bg-white">
            <Card>
              <CardHeader>
                <CardTitle>Recently Shortened URLs</CardTitle>
                <CardDescription>Your last 5 shortened URLs</CardDescription>
              </CardHeader>
              <CardContent>
                <ul className="space-y-4">
                  {recentUrls.map((url, index) => (
                    <li key={index} className="flex items-center justify-between">
                      <div className="flex-1 truncate mr-4">
                        <p className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{url.short}</p>
                        <p className="text-sm text-gray-500 dark:text-gray-400 truncate">{url.long}</p>
                      </div>
                      <Button size="icon" variant="outline" onClick={() => copyToClipboard(url.short)}>
                        <Clipboard className="h-4 w-4" />
                      </Button>
                      <Button size="icon" variant="outline" asChild>
                        <a href={url.short} target="_blank" rel="noopener noreferrer">
                          <ExternalLink className="h-4 w-4" />
                        </a>
                      </Button>
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </main>
    </div>
  )
}