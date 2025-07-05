package utils

import (
	"regexp"
	"strings"
)

// ChunkMessage splits a message into multiple chunks using ----- separator
// and falls back to sentence splitting for chunks that are too long
func ChunkMessage(response string, maxChunkLength int) []string {
	if response == "" {
		return []string{}
	}

	// Primary: Split on ----- separator
	chunks := strings.Split(response, "-----")
	
	var finalChunks []string
	
	for _, chunk := range chunks {
		cleanChunk := strings.TrimSpace(chunk)
		if cleanChunk == "" {
			continue
		}
		
		// If chunk is within limit, add it directly
		if len(cleanChunk) <= maxChunkLength {
			finalChunks = append(finalChunks, cleanChunk)
			continue
		}
		
		// If chunk is too long, split it further
		subChunks := splitLongChunk(cleanChunk, maxChunkLength)
		finalChunks = append(finalChunks, subChunks...)
	}
	
	// If no chunks were created (no ----- separators), treat entire response as one chunk
	if len(finalChunks) == 0 {
		if len(response) <= maxChunkLength {
			finalChunks = append(finalChunks, strings.TrimSpace(response))
		} else {
			finalChunks = splitLongChunk(strings.TrimSpace(response), maxChunkLength)
		}
	}
	
	return finalChunks
}

// splitLongChunk splits a long chunk into smaller pieces using sentence boundaries
func splitLongChunk(chunk string, maxLength int) []string {
	if len(chunk) <= maxLength {
		return []string{chunk}
	}
	
	var chunks []string
	
	// Try to split on sentences (. ! ?)
	sentences := splitOnSentences(chunk)
	
	currentChunk := ""
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}
		
		// If adding this sentence would exceed the limit
		if len(currentChunk)+len(sentence)+1 > maxLength {
			// Save current chunk if it has content
			if currentChunk != "" {
				chunks = append(chunks, strings.TrimSpace(currentChunk))
				currentChunk = ""
			}
			
			// If the sentence itself is too long, split it by words
			if len(sentence) > maxLength {
				wordChunks := splitByWords(sentence, maxLength)
				chunks = append(chunks, wordChunks...)
			} else {
				currentChunk = sentence
			}
		} else {
			// Add sentence to current chunk
			if currentChunk == "" {
				currentChunk = sentence
			} else {
				currentChunk += " " + sentence
			}
		}
	}
	
	// Add remaining chunk
	if currentChunk != "" {
		chunks = append(chunks, strings.TrimSpace(currentChunk))
	}
	
	return chunks
}

// splitOnSentences splits text into sentences using common sentence terminators
func splitOnSentences(text string) []string {
	// Regular expression to split on sentence boundaries
	// Looks for . ! ? followed by space and capital letter, or end of string
	re := regexp.MustCompile(`([.!?])\s+(?=[A-Z])|([.!?])$`)
	
	// Find all split positions
	indices := re.FindAllStringIndex(text, -1)
	
	if len(indices) == 0 {
		return []string{text}
	}
	
	var sentences []string
	lastEnd := 0
	
	for _, match := range indices {
		// Include the punctuation in the sentence
		end := match[1]
		sentence := text[lastEnd:end]
		if strings.TrimSpace(sentence) != "" {
			sentences = append(sentences, sentence)
		}
		lastEnd = end
	}
	
	// Add remaining text if any
	if lastEnd < len(text) {
		remaining := strings.TrimSpace(text[lastEnd:])
		if remaining != "" {
			sentences = append(sentences, remaining)
		}
	}
	
	return sentences
}

// splitByWords splits text by words when sentence splitting isn't enough
func splitByWords(text string, maxLength int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}
	
	var chunks []string
	currentChunk := ""
	
	for _, word := range words {
		// If adding this word would exceed the limit
		if len(currentChunk)+len(word)+1 > maxLength {
			// Save current chunk if it has content
			if currentChunk != "" {
				chunks = append(chunks, strings.TrimSpace(currentChunk))
				currentChunk = ""
			}
			
			// If the word itself is too long, just add it as its own chunk
			if len(word) > maxLength {
				chunks = append(chunks, word)
			} else {
				currentChunk = word
			}
		} else {
			// Add word to current chunk
			if currentChunk == "" {
				currentChunk = word
			} else {
				currentChunk += " " + word
			}
		}
	}
	
	// Add remaining chunk
	if currentChunk != "" {
		chunks = append(chunks, strings.TrimSpace(currentChunk))
	}
	
	return chunks
}

// CountWords counts the number of words in a text
func CountWords(text string) int {
	if text == "" {
		return 0
	}
	return len(strings.Fields(text))
}